{
  pkgs,
  system,
  self,
  ...
}: {
  scripts = {
    dx = {
      exec = ''$EDITOR "$(git rev-parse --show-toplevel)"/flake.nix'';
      description = "Edit flake.nix";
    };
    gx = {
      exec = ''$EDITOR "$(git rev-parse --show-toplevel)"/go.mod'';
      description = "Edit go.mod";
    };
    clean = {
      exec = ''git clean -fdx'';
      description = "Clean Project";
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
        templ generate "$REPO_ROOT"
        golangci-lint run "$REPO_ROOT"
        statix check "$REPO_ROOT"
        deadnix "$REPO_ROOT"/flake.nix
        nix flake check
      '';
      deps = with pkgs; [golangci-lint statix deadnix templ rustc cargo];
      description = "Run Nix/Go/Rust Linting Steps.";
    };
    generate-css = {
      exec = ''
        REPO_ROOT="$(git rev-parse --show-toplevel)"
        templ generate --log-level error "$REPO_ROOT"
        go run "$REPO_ROOT"/cmd/update-css --cwd "$REPO_ROOT"
        tailwindcss -i ./input.css \
            -o "$REPO_ROOT"/cmd/conneroh/_static/dist/style.css \
            --cwd "$REPO_ROOT"
      '';
      deps = with pkgs; [tailwindcss templ go];
      description = "Update the generated html and css files.";
    };
    generate-db = {
      exec = ''doppler run -- go run "$(git rev-parse --show-toplevel)"/cmd/update'';
      deps = with pkgs; [doppler];
      description = "Update the generated go files from the md docs.";
    };
    generate-reload = {
      exec = ''
        REPO_ROOT="$(git rev-parse --show-toplevel)"
        TEMPL_HASH=$(nix-hash --type sha256 --base32 "$REPO_ROOT"/cmd/conneroh/**/*.templ | sha256sum | cut -d' ' -f1)
        OLD_TEMPL_HASH=$(cat "$REPO_ROOT"/internal/cache/templ.hash)
        if [ "$OLD_TEMPL_HASH" != "$TEMPL_HASH" ]; then
          echo "OLD_TEMPL_HASH: $OLD_TEMPL_HASH; NEW_TEMPL_HASH: $TEMPL_HASH"
          generate-css
          echo "$TEMPL_HASH" > ./internal/cache/templ.hash
        fi
        DOCS_HASH=$(nix-hash --type sha256 --base32 ./internal/data/**/*.md | sha256sum | cut -d' ' -f1)
        OLD_DOCS_HASH=$(cat "$REPO_ROOT"/internal/cache/docs.hash)
        if [ "$OLD_DOCS_HASH" != "$DOCS_HASH" ]; then
          echo "OLD_DOCS_HASH: $OLD_DOCS_HASH; NEW_DOCS_HASH: $DOCS_HASH"
          generate-db
          echo "$DOCS_HASH" > ./internal/cache/docs.hash
        fi
      '';
      deps = with self.packages."${system}"; [generate-db generate-css];
      description = "Code Generation Steps for specific directory changes.";
    };
    generate-js = {
      exec = ''bun build "$(git rev-parse --show-toplevel)"/index.js --minify --minify-syntax --minify-whitespace  --minify-identifiers --outdir "$(git rev-parse --show-toplevel)"/cmd/conneroh/_static/dist/'';
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
      deps = with self.packages."${system}"; [generate-css generate-db generate-js];
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
      exec = ''go get -u "$(git rev-parse --show-toplevel)"/...'';
      deps = with pkgs; [go git];
      description = "Update go.sum";
    };
    generate-templates = {
      exec = ''templ generate "$(git rev-parse --show-toplevel)"'';
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
    pr-review = {
      exec = ''
        set -eo pipefail

        # Required dependencies check
        for cmd in jq flyctl cut; do
          if ! command -v "$cmd" >/dev/null 2>&1; then
            echo "Error: Required command '$cmd' not found. Ensure it's included in runtimeInputs."
            exit 1
          fi
        done

        # Handle working directory change
        if [ -n "''${INPUT_PATH}" ]; then
          cd "''${INPUT_PATH}" || exit 1
        fi

        if [ -z "''${GITHUB_EVENT_PATH}" ]; then
          if [ -n "$1" ] && [ -f "$1" ]; then
            GITHUB_EVENT_PATH="$1"
          else
            echo "Error: GITHUB_EVENT_PATH not set and no event file provided as argument"
            exit 1
          fi
        fi
        PR_NUMBER=$(jq -r '.pull_request.number // .number // empty' "''${GITHUB_EVENT_PATH}")
        if [ -z "''${PR_NUMBER}" ]; then
          echo "Error: Could not extract PR number. This action only supports pull_request events."
          exit 1
        fi

        # Ensure repository info is available - use cut instead of parameter expansion
        if [ -n "''${GITHUB_REPOSITORY}" ]; then
          # Extract owner and name using cut rather than parameter expansion
          GITHUB_REPOSITORY_OWNER=$(echo "''${GITHUB_REPOSITORY}" | cut -d '/' -f 1)
          GITHUB_REPOSITORY_NAME=$(echo "''${GITHUB_REPOSITORY}" | cut -d '/' -f 2)
        fi

        # If we don't have the repository information, we need both parts
        if [ -z "''${GITHUB_REPOSITORY_OWNER}" ] || [ -z "''${GITHUB_REPOSITORY_NAME}" ]; then
          echo "Error: Repository owner or name not available. Please set GITHUB_REPOSITORY_OWNER and GITHUB_REPOSITORY_NAME."
          exit 1
        fi

        # Extract event type
        EVENT_TYPE=$(jq -r '.action // empty' "''${GITHUB_EVENT_PATH}")

        # Prepare app name with PR number, ensuring it follows Fly.io naming conventions
        app="''${INPUT_NAME:-pr-''${PR_NUMBER}-''${GITHUB_REPOSITORY_OWNER}-''${GITHUB_REPOSITORY_NAME}}"
        app=$(echo "''${app}" | tr '_' '-')  # Change underscores to hyphens

        # Configuration
        region="''${INPUT_REGION:-''${FLY_REGION:-iad}}"
        org="''${INPUT_ORG:-''${FLY_ORG:-personal}}"
        image="''${INPUT_IMAGE}"
        config="''${INPUT_CONFIG:-fly.toml}"
        build_args=""
        build_secrets=""
        output_file="''${GITHUB_OUTPUT:-fly-output.env}"

        # Set default values for optional inputs
        : "''${INPUT_HA:=false}"
        : "''${INPUT_VMSIZE:=shared-cpu-1x}"
        : "''${INPUT_CPUKIND:=shared}"
        : "''${INPUT_CPU:=1}"
        : "''${INPUT_MEMORY:=256MB}"

        # Safety check: ensure app name contains PR number
        if ! echo "''${app}" | grep -q "''${PR_NUMBER}"; then
          echo "For safety, this action requires the app's name to contain the PR number."
          exit 1
        fi

        # Handle PR closure - remove the Fly app if it exists
        if [ "''${EVENT_TYPE}" = "closed" ]; then
          echo "PR was closed. Attempting to destroy the Fly app..."
          flyctl apps destroy "''${app}" -y || true
          exit 0
        fi

        # Process build arguments
        if [ -n "''${INPUT_BUILD_ARGS}" ]; then
          for ARG in $(echo "''${INPUT_BUILD_ARGS}" | tr " " "\n"); do
            build_args="$build_args --build-arg ''${ARG}"
          done
        fi

        # Process build secrets
        if [ -n "''${INPUT_BUILD_SECRETS}" ]; then
          for ARG in $(echo "''${INPUT_BUILD_SECRETS}" | tr " " "\n"); do
            build_secrets="$build_secrets --build-secret ''${ARG}"
          done
        fi

        # Deploy the Fly app, creating it first if needed
        if ! flyctl status --app "''${app}" 2>/dev/null; then
          echo "App doesn't exist yet. Creating new app: ''${app}"
          # Backup the original config file since 'flyctl launch' can modify it
          cp "''${config}" "''${config}.bak"
          flyctl launch --no-deploy --copy-config --name "''${app}" --image "''${image}" \
            --region "''${region}" --org "''${org}"
          # Restore the original config file
          cp "''${config}.bak" "''${config}"
        fi

        # Set secrets if provided
        if [ -n "''${INPUT_SECRETS}" ]; then
          echo "''${INPUT_SECRETS}" | tr " " "\n" | flyctl secrets import --app "''${app}"
        fi

        # Attach postgres cluster if specified
        if [ -n "''${INPUT_POSTGRES}" ]; then
          echo "Attaching Postgres cluster: ''${INPUT_POSTGRES}"
          flyctl postgres attach "''${INPUT_POSTGRES}" --app "''${app}" || true
        fi

        # Deploy the application
        echo "Contents of config ''${config} file: " && cat "''${config}"

        if [ -n "''${INPUT_VM}" ]; then
          flyctl deploy --config "''${config}" --app "''${app}" --regions "''${region}" \
            --image "''${image}" --strategy immediate --ha="''${INPUT_HA}" \
            --vm-size "''${INPUT_VMSIZE}"
        else
          flyctl deploy --config "''${config}" --app "''${app}" --regions "''${region}" \
            --image "''${image}" --strategy immediate --ha="''${INPUT_HA}" \
            --vm-cpu-kind "''${INPUT_CPUKIND}" --vm-cpus "''${INPUT_CPU}" \
            --vm-memory "''${INPUT_MEMORY}"
        fi

        # Get status and write output
        flyctl status --app "''${app}" --json > status.json
        hostname=$(jq -r '.Hostname // empty' status.json)
        appid=$(jq -r '.ID // empty' status.json)

        # Write outputs in a way that works both in GitHub Actions and standalone
        if [ -n "''${GITHUB_OUTPUT}" ]; then
          # GitHub Actions environment - use a block instead of individual redirects
          {
            echo "hostname=''${hostname}"
            echo "url=https://''${hostname}"
            echo "id=''${appid}"
            echo "name=''${app}"
          } >> "''${GITHUB_OUTPUT}"
        else
          # Standalone environment - write to file and stdout
          {
            echo "hostname=''${hostname}"
            echo "url=https://''${hostname}"
            echo "id=''${appid}"
            echo "name=''${app}"
          } >> "''${output_file}"

          echo "Deployment completed successfully!"
          echo "App: ''${app}"
          echo "URL: https://''${hostname}"
          echo "App ID: ''${appid}"
          echo "Output written to: ''${output_file}"
        fi
      '';
      deps = with pkgs; [];
      description = "Create a PR Review Deployment";
    };
  };
}
