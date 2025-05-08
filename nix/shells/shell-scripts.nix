{
  pkgs,
  self',
  ...
}: rec {
  scripts = {
    dx = {
      exec = ''
        REPO_ROOT="$(git rev-parse --show-toplevel)"
        $EDITOR "$REPO_ROOT"/flake.nix
      '';
      deps = [pkgs.git];
      description = "Edit flake.nix";
    };
    gx = {
      exec = ''
        REPO_ROOT="$(git rev-parse --show-toplevel)"
        $EDITOR "$REPO_ROOT"/go.mod
      '';
      deps = [pkgs.git];
      description = "Edit go.mod";
    };
    clean = {
      exec = ''git clean -fdx'';
      description = "Clean Project";
      deps = [pkgs.git];
    };
    reset-db = {
      exec = ''
        for f in "./master.db" "./master.db-shm" "./master.db-wal"; do
          rm -f "$f"
        done
      '';
      description = "Reset the database";
    };
    tests = {
      exec = ''
        REPO_ROOT="$(git rev-parse --show-toplevel)"
        go test -v "$REPO_ROOT"/...
      '';
      deps = [pkgs.go];
      description = "Run all go tests";
    };
    lint = {
      exec = ''
        REPO_ROOT="$(git rev-parse --show-toplevel)"
        templ generate

        golangci-lint run
        statix check "$REPO_ROOT"/flake.nix
        deadnix "$REPO_ROOT"/flake.nix
      '';
      deps = with pkgs; [golangci-lint statix deadnix templ];
      description = "Run Nix/Go Linting Steps.";
    };
    generate-css = {
      exec = ''
        REPO_ROOT="$(git rev-parse --show-toplevel)"

        templ generate --log-level error
        go run "$REPO_ROOT"/cmd/update-css --cwd "$REPO_ROOT"
        tailwindcss -i ./input.css \
            -o "$REPO_ROOT"/cmd/conneroh/_static/dist/style.css \
            --cwd "$REPO_ROOT"
      '';
      deps = with pkgs; [tailwindcss templ go];
      description = "Update the generated html and css files.";
    };
    generate-db = {
      exec = ''
        # Run update from PATH if available, or build it if needed
        if command -v update >/dev/null 2>&1; then
          doppler run -- update
        else
          echo "update command not found, building it first..."
          nix build .#update --no-link
          nix run .#update
        fi
      '';
      deps = with pkgs; [doppler];
      description = "Update the generated go files from the md docs.";
    };
    generate-reload = {
      exec = ''
        REPO_ROOT="$(git rev-parse --show-toplevel)" # needed

        TEMPL_HASH=$(nix-hash --type sha256 --base32 "$REPO_ROOT"/cmd/conneroh/**/*.templ | sha256sum | cut -d' ' -f1)
        OLD_TEMPL_HASH=$(cat "$REPO_ROOT"/internal/cache/templ.hash)

        if [ "$OLD_TEMPL_HASH" != "$TEMPL_HASH" ]; then
          echo "OLD_TEMPL_HASH: $OLD_TEMPL_HASH; NEW_TEMPL_HASH: $TEMPL_HASH"
          generate-css
          echo "$TEMPL_HASH" > ./internal/cache/templ.hash
        fi

        DOCS_HASH=$(nix-hash --type sha256 --base32 ./internal/data/docs/ | sha256sum | cut -d' ' -f1)
        OLD_DOCS_HASH=$(cat "$REPO_ROOT"/internal/cache/docs.hash)

        if [ "$OLD_DOCS_HASH" != "$DOCS_HASH" ]; then
          echo "OLD_DOCS_HASH: $OLD_DOCS_HASH; NEW_DOCS_HASH: $DOCS_HASH"
          generate-db
          echo "$DOCS_HASH" > ./internal/cache/docs.hash
        fi
      '';
      deps = with self'.packages; [generate-db generate-css];
      description = "Code Generation Steps for specific directory changes.";
    };
    generate-js = {
      exec = ''
        REPO_ROOT="$(git rev-parse --show-toplevel)"
        bun build \
              "$REPO_ROOT"/index.js \
              --minify \
              --minify-syntax \
              --minify-whitespace  \
              --minify-identifiers \
              --outdir "$REPO_ROOT"/cmd/conneroh/_static/dist/
      '';
      deps = with pkgs; [bun git];
      description = "Generate JS files";
    };
    generate-all = {
      exec = ''
        generate-css &
        generate-db &
        generate-js &
        wait
      '';
      deps = with self'.packages; [generate-css generate-db generate-js];
      description = "Generate all files in parallel";
    };
    format = {
      exec = ''
        cd "$(git rev-parse --show-toplevel)"
        go fmt ./...
        git ls-files \
            --others \
            --exclude-standard \
            --cached \
            -- '*.js' '*.ts' '*.css' '*.md' '*.json' \
            | xargs prettier --write
        golines \
            -l \
            -w \
            --max-len=80 \
            --shorten-comments \
            --ignored-dirs=.direnv .
        cd -
      '';
      deps = with pkgs; [go git golines];
      description = "Format code files";
    };
    goSumUpdate = {
      exec = ''
        echo "Updating go.sum..."
        go get -u ./...
      '';
      deps = with pkgs; [go git];
      description = "Update go.sum";
    };
    generate-templates = {
      exec = ''
        templ generate
      '';
      deps = with pkgs; [templ];
      description = "Generate templates";
    };
    run = {
      exec = ''
        export DEBUG=true
        cd "$(git rev-parse --show-toplevel)" && air
      '';
      deps = with pkgs; [air git];
      description = "Run the application with air for hot reloading";
    };
    live-ci = {
      exec = ''
        go run ./cmd/live-ci/main.go
      '';
      env = {
        DEBUG = "true";
      };
      deps = with pkgs;
        [
          playwright-driver # Browser Archives and Driver Scripts
          nodejs_20 # Required for Playwright driver
          pkg-config # Needed for some browser dependencies
          at-spi2-core # Accessibility support
          cairo # 2D graphics library
          cups # Printing system
          dbus # Message bus system
          expat # XML parser
          ffmpeg # Media processing
          fontconfig # Font configuration and customization
          freetype # Font rendering engine
          gdk-pixbuf # Image loading library
          glib # Low-level core library
          gtk3 # GUI toolkit
          go
        ]
        ++ (with pkgs;
          lib.optionals stdenv.isDarwin [
            libiconv
          ])
        ++ (with pkgs;
          lib.optionals stdenv.isLinux [
            chromium # Chromium browser
            xorg.libXcomposite # X11 Composite extension - needed by browsers
            xorg.libXdamage # X11 Damage extension - needed by browsers
            xorg.libXfixes # X11 Fixes extension - needed by browsers
            xorg.libXrandr # X11 RandR extension - needed by browsers
            xorg.libX11 # X11 client-side library
            xorg.libxcb # X11 C Bindings library
            mesa # OpenGL implementation
            alsa-lib # Audio library
            nss # Network Security Services
            nspr # NetScape Portable Runtime
            pango # Text layout and rendering
          ]);
      description = "Run the application with air for hot reloading";
    };
  };
  scriptPackages =
    pkgs.lib.mapAttrs
    (
      name: script:
      # Create a script with dependencies
        pkgs.writeShellApplication {
          inherit name;
          text = script.exec;
          # Add runtime dependencies
          runtimeInputs = script.deps or [];
        }
    )
    scripts;
}
