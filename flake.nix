{
  description = "Personal Website for Conner Ohnesorge";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";

    flake-utils = {
      url = "github:numtide/flake-utils";
      inputs.systems.follows = "systems";
    };

    nix2container = {
      url = "github:nlewo/nix2container";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.flake-utils.follows = "flake-utils";
    };

    twerge = {
      url = "github:conneroisu/twerge?tag=v0.2.3";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.flake-utils.follows = "flake-utils";
    };

    mk-shell-bin.url = "github:rrbutani/nix-mk-shell-bin";

    bun2nix.url = "github:baileyluTCD/bun2nix";

    systems.url = "github:nix-systems/default";
  };

  nixConfig = {
    extra-substituters = ''
      https://cache.nixos.org
      https://nix-community.cachix.org
      https://devenv.cachix.org
      https://conneroisu.cachix.org
    '';
    extra-trusted-public-keys = ''
      cache.nixos.org-1:6NCHdD59X431o0gWypbMrAURkbJ16ZPMQFGspcDShjY=
      nix-community.cachix.org-1:mB9FSh9qf2dCimDSUo8Zy7bkq5CX+/rkCWyvRCYg3Fs=
      devenv.cachix.org-1:w1cLUi8dv3hnoSPGAuibQv+f9TZLr6cv/Hm9XgU50cw=
      conneroisu.cachix.org-1:PgOlJ8/5i/XBz2HhKZIYBSxNiyzalr1B/63T74lRcU0=
    '';
    extra-experimental-features = "nix-command flakes";
  };

  outputs = inputs @ {
    self,
    flake-utils,
    ...
  }:
    flake-utils.lib.eachSystem [
      "x86_64-linux"
      "i686-linux"
      "x86_64-darwin"
      "aarch64-linux"
      "aarch64-darwin"
    ] (system: let
      overlay = final: prev: {final.go = prev.go_1_24;};
      pkgs = import inputs.nixpkgs {
        inherit system;
        overlays = [
          overlay
          inputs.twerge.overlays."${system}".default
        ];
      };
      buildGoModule = pkgs.buildGoModule.override {go = pkgs.go_1_24;};
      buildWithSpecificGo = pkg: pkg.override {inherit buildGoModule;};

      scripts = {
        dx = {
          exec = ''$EDITOR $REPO_ROOT/flake.nix'';
          description = "Edit flake.nix";
        };
        gx = {
          exec = ''$EDITOR $REPO_ROOT/go.mod'';
          description = "Edit go.mod";
        };
        clean = {
          exec = ''${pkgs.git}/bin/git clean -fdx'';
          description = "Clean Project";
        };
        tests = {
          exec = ''${pkgs.go}/bin/go test -v ./...'';
          description = "Run all go tests";
        };
        lint = {
          exec = ''
            ${pkgs.golangci-lint}/bin/golangci-lint run
            ${pkgs.statix}/bin/statix check $REPO_ROOT/flake.nix
            ${pkgs.deadnix}/bin/deadnix $REPO_ROOT/flake.nix
          '';
          description = "Run Linting Steps.";
        };
        build = {
          exec = ''
            nix build --accept-flake-config .#packages.x86_64-linux.conneroh
          '';
          description = "Build the package";
        };
        update = {
          exec = ''
            ${pkgs.doppler}/bin/doppler run -- ${pkgs.go}/bin/go run $REPO_ROOT/cmd/update --cwd $REPO_ROOT
          '';
          description = "Update the generated go files.";
        };
        generate-reload = {
          exec = ''
            export REPO_ROOT=$(git rev-parse --show-toplevel) # needed
            cd $REPO_ROOT
            if ${pkgs.hasher}/bin/hasher -dir "$REPO_ROOT/cmd/conneroh/views" -v -exclude "*_templ.go"; then
              echo ""
              if ${pkgs.hasher}/bin/hasher -dir "$REPO_ROOT/internal/data/docs" -v -exclude "*_templ.go"; then
                echo ""
              else
                echo "Changes detected in docs, running update script..."
                ${pkgs.doppler}/bin/doppler run -- ${pkgs.go}/bin/go run $REPO_ROOT/cmd/update --cwd $REPO_ROOT
              fi
            else
              echo "Changes detected in templates, running update script..."
              ${pkgs.templ}/bin/templ generate --log-level error
              go run $REPO_ROOT/cmd/update-css --cwd $REPO_ROOT
              ${pkgs.templ}/bin/templ generate --log-level error
              ${pkgs.tailwindcss}/bin/tailwindcss -m -i ./input.css -o ./cmd/conneroh/_static/dist/style.css --cwd $REPO_ROOT
            fi
            cd -
          '';
          description = "Code Generation Steps for specific directory changes.";
        };
        generate-all = {
          exec = ''
            ${pkgs.templ}/bin/templ generate
            ${pkgs.go}/bin/go run $REPO_ROOT/cmd/update-css --cwd $REPO_ROOT
            ${pkgs.tailwindcss}/bin/tailwindcss \
                --minify \
                -i ./input.css \
                -o ./cmd/conneroh/_static/dist/style.css \
                --cwd .
          '';
          description = "Generate all files in parallel";
        };
        format = {
          exec = ''
            cd $(git rev-parse --show-toplevel)

            ${pkgs.go}/bin/go fmt ./...

            ${pkgs.git}/bin/git ls-files \
              --others \
              --exclude-standard \
              --cached \
              -- '*.js' '*.ts' '*.css' '*.md' '*.json' \
              | xargs prettier --write

            ${pkgs.golines}/bin/golines \
              -l \
              -w \
              --max-len=80 \
              --shorten-comments \
              --ignored-dirs=.direnv .

            cd -
          '';
          description = "Format code files";
        };
        run = {
          exec = ''cd $REPO_ROOT && air'';
          description = "Run the application with air for hot reloading";
        };
      };

      # Convert scripts to packages
      scriptPackages =
        pkgs.lib.mapAttrsToList
        (name: script: pkgs.writeShellScriptBin name script.exec)
        scripts;
    in rec {
      devShells.default = pkgs.mkShell {
        shellHook = ''
          export REPO_ROOT=$(git rev-parse --show-toplevel)
          export CGO_CFLAGS="-O2"

          export PLAYWRIGHT_BROWSERS_PATH=${pkgs.playwright-driver.browsers}
          export PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD=1
          export PLAYWRIGHT_NODEJS_PATH=${pkgs.nodejs_20}/bin/node

          # Browser executable paths
          export PLAYWRIGHT_CHROMIUM_EXECUTABLE_PATH=${"${pkgs.playwright-driver.browsers}/chromium-1155"}
          export PLAYWRIGHT_FIREFOX_EXECUTABLE_PATH=${"${pkgs.playwright-driver.browsers}/firefox-1471"}
          export PLAYWRIGHT_WEBKIT_EXECUTABLE_PATH=${"${pkgs.playwright-driver.browsers}/webkit-2123"}

          echo "Playwright configured with:"
          echo "  - Browsers directory: $PLAYWRIGHT_BROWSERS_PATH"
          echo "  - Node.js path: $PLAYWRIGHT_NODEJS_PATH"
          echo "  - Chromium path: $PLAYWRIGHT_CHROMIUM_EXECUTABLE_PATH"
          echo "  - Firefox path: $PLAYWRIGHT_FIREFOX_EXECUTABLE_PATH"
          echo "  - WebKit path: $PLAYWRIGHT_WEBKIT_EXECUTABLE_PATH"

          # Print available commands
          echo "Available commands:"
          ${pkgs.lib.concatStringsSep "\n" (
            pkgs.lib.mapAttrsToList (
              name: script: ''echo "  ${name} - ${script.description}"''
            )
            scripts
          )}
        '';
        packages = with pkgs;
          [
            # Nix
            alejandra
            nixd
            statix
            deadnix
            inputs.bun2nix.defaultPackage.${pkgs.system}.bin

            # Go Tools
            go_1_24
            air
            templ
            pprof
            revive
            golangci-lint
            (buildWithSpecificGo gopls)
            (buildWithSpecificGo templ)
            (buildWithSpecificGo golines)
            (buildWithSpecificGo golangci-lint-langserver)
            (buildWithSpecificGo gomarkdoc)
            (buildWithSpecificGo gotests)
            (buildWithSpecificGo gotools)
            (buildWithSpecificGo reftools)
            graphviz

            # Web
            tailwindcss
            tailwindcss-language-server
            bun
            nodePackages.typescript-language-server
            nodePackages.prettier

            # Infra
            flyctl
            wireguard-tools
            openssl.dev
            skopeo
            inputs.twerge.packages."${pkgs.system}".hasher

            # Playwright

            playwright-driver # Provides browser archives and driver scripts
            (
              if pkgs.stdenv.isDarwin
              then pkgs.darwin.apple_sdk.frameworks.WebKit
              else pkgs.webkitgtk
            ) # WebKit browser
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
          ]
          ++ pkgs.lib.optionals pkgs.stdenv.isDarwin [
            # macOS-specific dependencies
            libiconv
          ]
          ++ pkgs.lib.optionals pkgs.stdenv.isLinux [
            # Linux-specific dependencies
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
          ]
          # Add the generated script packages
          ++ scriptPackages;
      };

      packages = let
        src = ./.;
        name = "conneroh.com";
        fly-name = "conneroh-com";
        fly-name-dev = "conneroh-com-dev";
        vendorHash = "sha256-K52okJQZ/y1VQb8ob4zcbuNC8hjhgUTPDBeVA1FJCKA=";
        created = "now";
        tag = "latest";
        version = self.shortRev or "dirty";
        nativeBuildInputs = [];
        preBuild = ''
          ${pkgs.templ}/bin/templ generate
          ${pkgs.tailwindcss}/bin/tailwindcss \
              --minify \
              -i ./input.css \
              -o ./cmd/conneroh/_static/dist/style.css \
              --cwd .
        '';
        config = {
          WorkingDir = "/root";
          Cmd = ["/bin/conneroh.com"];
          ExposedPorts = {
            "8080/tcp" = {};
          };
          Env = [
            "SSL_CERT_FILE=${pkgs.cacert}/etc/ssl/certs/ca-bundle.crt"
            "NIX_SSL_CERT_FILE=${pkgs.cacert}/etc/ssl/certs/ca-bundle.crt"
          ];
        };
      in rec {
        conneroh = buildGoModule {
          inherit vendorHash name src preBuild nativeBuildInputs version;
          subPackages = ["."];
        };
        C-conneroh = pkgs.dockerTools.buildLayeredImage {
          inherit name config tag created;
          contents = [
            conneroh
            pkgs.cacert
          ];
          extraCommands = ''
            echo "$(git rev-parse HEAD)" > REVISION
          '';
        };
        C-conneroh-dev = pkgs.dockerTools.buildLayeredImage {
          inherit name config tag created;
          contents = [
            conneroh
            pkgs.cacert
          ];
          extraCommands = ''
            echo "$(git rev-parse HEAD)" > REVISION
          '';
        };
        deployPackage = pkgs.writeShellScriptBin "deploy" ''
          set -e

          if [ -z "$FLY_AUTH_TOKEN" ]; then
            echo "FLY_AUTH_TOKEN is not set. Getting it from doppler..."
            FLY_AUTH_TOKEN=$(${pkgs.doppler}/bin/doppler secrets get --plain FLY_AUTH_TOKEN)
          fi

          echo "Copying image to Fly.io registry..."
          ${pkgs.skopeo}/bin/skopeo copy \
            --insecure-policy \
            docker-archive:"${C-conneroh}" \
            docker://registry.fly.io/${fly-name}:latest \
            --dest-creds x:"$FLY_AUTH_TOKEN" \
            --format v2s2

          echo "Deploying to Fly.io..."
          ${pkgs.flyctl}/bin/fly deploy \
            --remote-only \
            -c ${./fly.toml} \
            -i registry.fly.io/${fly-name} \
            -t "$FLY_AUTH_TOKEN"
        '';

        deployPackageDev = pkgs.writeShellScriptBin "deploy-package-dev" ''
          set -e

          if [ -z "$FLY_DEV_AUTH_TOKEN" ]; then
            echo "FLY_AUTH_TOKEN is not set. Getting it from doppler..."
            FLY_DEV_AUTH_TOKEN=$(${pkgs.doppler}/bin/doppler secrets get --plain FLY_DEV_AUTH_TOKEN)
          fi

          echo "Copying dev image to Fly.io registry..."
          ${pkgs.skopeo}/bin/skopeo copy \
            --insecure-policy \
            docker-archive:"${C-conneroh-dev}" \
            docker://registry.fly.io/${fly-name-dev}:latest \
            --dest-creds x:"$FLY_DEV_AUTH_TOKEN" \
            --format v2s2

          echo "Deploying to Fly.io..."
          ${pkgs.flyctl}/bin/fly deploy \
            --remote-only \
            -c ${./fly.dev.toml} \
            -i registry.fly.io/${fly-name-dev}:latest \
            -t "$FLY_DEV_AUTH_TOKEN"
        '';
      };
    });
}
