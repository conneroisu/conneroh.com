{
  description = "Personal Website for Conner Ohnesorge";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    systems.url = "github:nix-systems/default";
    flake-utils = {
      url = "github:numtide/flake-utils";
      inputs.systems.follows = "systems";
    };
    bun2nix.url = "github:baileyluTCD/bun2nix";
  };

  nixConfig = {
    extra-substituters = ''https://conneroisu.cachix.org'';
    extra-trusted-public-keys = ''conneroisu.cachix.org-1:PgOlJ8/5i/XBz2HhKZIYBSxNiyzalr1B/63T74lRcU0='';
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
        ];
        config.allowUnfree = true;
      };
      buildWithSpecificGo = pkg: pkg.override {buildGoModule = pkgs.buildGo124Module;};
    in rec {
      devShells.default = let
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
          update = {
            exec = ''
              ${pkgs.doppler}/bin/doppler run -- ${packages.update}/bin/update -cwd $REPO_ROOT -jobs 20
            '';
            description = "Update the generated go files.";
          };
          generate-reload = {
            exec = ''
              export REPO_ROOT=$(git rev-parse --show-toplevel) # needed
              function gen_css() {
                ${pkgs.templ}/bin/templ generate --log-level error
                ${pkgs.go}/bin/go run $REPO_ROOT/cmd/update-css --cwd $REPO_ROOT
                ${pkgs.tailwindcss}/bin/tailwindcss -m -i ./input.css -o $REPO_ROOT/cmd/conneroh/_static/dist/style.css --cwd $REPO_ROOT
              }
              function gen_doc() {
                ${pkgs.doppler}/bin/doppler run -- ${packages.update}/bin/update -cwd $REPO_ROOT -jobs 20
              }

              TEMPL_HASH=$(nix-hash --type sha256 --base32 $REPO_ROOT/cmd/conneroh/**/*.templ | sha256sum | cut -d' ' -f1)
              echo "TEMPL_HASH: $TEMPL_HASH"
              OLD_TEMPL_HASH=$(cat $REPO_ROOT/cmd/conneroh/_static/dist/templ.hash)
              if [ "$OLD_TEMPL_HASH" != "$TEMPL_HASH" ]; then
                echo "templ change"
                gen_css
                echo "$TEMPL_HASH" > ./internal/cache/templ.hash
              fi

              DOCS_HASH=$(nix-hash --type sha256 --base32 ./internal/data/docs/ | sha256sum | cut -d' ' -f1)
              echo "DOCS_HASH: $DOCS_HASH"
              OLD_DOCS_HASH=$(cat $REPO_ROOT/cmd/conneroh/_static/dist/docs.hash)
              if [ "$OLD_DOCS_HASH" != "$DOCS_HASH" ]; then
                echo "docs change"
                gen_doc
                echo "$DOCS_HASH" > ./internal/cache/docs.hash
              fi
            '';
            description = "Code Generation Steps for specific directory changes.";
          };
          generate-all = {
            exec = ''
              export REPO_ROOT=$(git rev-parse --show-toplevel)
              ${pkgs.templ}/bin/templ generate
              ${pkgs.go}/bin/go run $REPO_ROOT/cmd/update-css --cwd $REPO_ROOT
              ${pkgs.tailwindcss}/bin/tailwindcss \
                  --minify \
                  -i ./input.css \
                  -o ./cmd/conneroh/_static/dist/style.css \
                  --cwd $REPO_ROOT
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
          generate-js = {
            exec = ''
              ${pkgs.bun}/bin/bun build \
                  $REPO_ROOT/index.js \
                  --minify \
                  --minify-syntax \
                  --minify-whitespace  \
                  --minify-identifiers \
                  --outdir $REPO_ROOT/cmd/conneroh/_static/dist/ &
            '';
            description = "Generate JS files";
          };
          run = {
            exec = "cd $REPO_ROOT && air";
            description = "Run the application with air for hot reloading";
          };
        };

        scriptPackages =
          pkgs.lib.mapAttrsToList
          (name: script: pkgs.writeShellScriptBin name script.exec)
          scripts;
      in
        pkgs.mkShell {
          shellHook = ''
            export REPO_ROOT=$(git rev-parse --show-toplevel)
            export CGO_CFLAGS="-O2"

            export PLAYWRIGHT_BROWSERS_PATH=${pkgs.playwright-driver.browsers}
            export PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD=1
            export PLAYWRIGHT_NODEJS_PATH=${pkgs.nodejs_20}/bin/node

            # Browser executable paths
            export PLAYWRIGHT_CHROMIUM_EXECUTABLE_PATH=${"${pkgs.playwright-driver.browsers}/chromium-1155"}

            echo "Playwright configured with:"
            echo "  - Browsers directory: $PLAYWRIGHT_BROWSERS_PATH"
            echo "  - Node.js path: $PLAYWRIGHT_NODEJS_PATH"
            echo "  - Chromium path: $PLAYWRIGHT_CHROMIUM_EXECUTABLE_PATH"

            # Print available commands
            echo "Available commands:"
            ${pkgs.lib.concatStringsSep "\n" (
              pkgs.lib.mapAttrsToList (
                name: script: ''echo "  ${name} - ${script.description}"''
              )
              scripts
            )}

            echo "Git Status:"
            ${pkgs.git}/bin/git status
          '';
          packages = with pkgs;
            [
              alejandra # Nix
              nixd
              statix
              deadnix

              go_1_24 # Go Tools
              air
              templ
              golangci-lint
              (buildWithSpecificGo revive)
              (buildWithSpecificGo gopls)
              (buildWithSpecificGo templ)
              (buildWithSpecificGo golines)
              (buildWithSpecificGo golangci-lint-langserver)
              (buildWithSpecificGo gomarkdoc)
              (buildWithSpecificGo gotests)
              (buildWithSpecificGo gotools)
              (buildWithSpecificGo reftools)
              pprof
              graphviz

              tailwindcss # Web
              tailwindcss-language-server
              bun
              nodePackages.typescript-language-server
              nodePackages.prettier
              inputs.bun2nix.defaultPackage.${pkgs.system}.bin

              flyctl # Infra
              openssl.dev
              skopeo
              consul

              playwright-driver # Browser Archives and Driver Scripts
              (
                if stdenv.isDarwin
                then darwin.apple_sdk.frameworks.WebKit
                else webkitgtk
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
              ])
            ++ scriptPackages;
        };

      packages = let
        internal_port = 8080;
        force_https = true;
        processes = ["app"];
        src = ./.;
        vendorHash = null;
      in rec {
        conneroh = pkgs.buildGo124Module {
          inherit src vendorHash;
          name = "conneroh.com";
          version = self.shortRev or "dirty";
          preBuild = ''
            ${pkgs.templ}/bin/templ generate
            ${pkgs.tailwindcss}/bin/tailwindcss \
                --minify \
                -i ./input.css \
                -o ./cmd/conneroh/_static/dist/style.css \
                --cwd .
          '';
          subPackages = ["."];
        };
        update = pkgs.buildGoModule {
          inherit src vendorHash;
          version = self.shortRev or "dirty";
          name = "update";
          subPackages = ["./cmd/update"];
        };
        C-conneroh = pkgs.dockerTools.buildImage {
          created = "now";
          tag = "latest";
          name = "conneroh";
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
          copyToRoot = [conneroh pkgs.cacert];
        };
        deployPackage = let
          settingsFormat = pkgs.formats.toml {};

          flyDevConfig = {
            app = "conneroh-com-dev";
            primary_region = "ord";
            build = {};
            http_service = {
              inherit internal_port force_https processes;
              auto_stop_machines = "stop";
              auto_start_machines = true;
              min_machines_running = 0;
            };
            vm = [
              {
                memory = "512M";
                cpu_kind = "shared";
                cpus = 1;
              }
            ];
          };

          flyProdConfig = {
            app = "conneroh-com";
            primary_region = "ord";
            build = {};
            http_service = {
              inherit internal_port force_https processes;
              auto_stop_machines = "stop";
              auto_start_machines = true;
              min_machines_running = 0;
            };
            vm = [
              {
                memory = "1gb";
                cpu_kind = "shared";
                cpus = 2;
              }
            ];
          };

          flyDevToml = settingsFormat.generate "fly.dev.toml" flyDevConfig;
          flyProdToml = settingsFormat.generate "fly.toml" flyProdConfig;
        in
          pkgs.writeShellScriptBin "deploy" ''
            set -e
            arg=$1
            TOKEN=""
            FLY_NAME=""
            CONFIG_FILE=""

            [ -z "$arg" ] && { echo "No argument provided. Please provide a target environment. (dev or prod)"; exit 1; }

            if [ "$arg" = "dev" ]; then
              [ -z "$FLY_DEV_AUTH_TOKEN" ] && FLY_DEV_AUTH_TOKEN=$(doppler secrets get --plain FLY_DEV_AUTH_TOKEN)
              TOKEN="$FLY_DEV_AUTH_TOKEN"
              export FLY_NAME="conneroh-com-dev"
              export CONFIG_FILE=${flyDevToml}
            else
              [ -z "$FLY_AUTH_TOKEN" ] && FLY_AUTH_TOKEN=$(doppler secrets get --plain FLY_AUTH_TOKEN)
              TOKEN="$FLY_AUTH_TOKEN"
              export FLY_NAME="conneroh-com"
              export CONFIG_FILE=${flyProdToml}
            fi

            REGISTY="registry.fly.io/$FLY_NAME"
            echo "Copying image to Fly.io... to $REGISTY"

            ${pkgs.skopeo}/bin/skopeo copy \
              --insecure-policy \
              docker-archive:"${C-conneroh}" \
              "docker://$REGISTY:latest" \
              --dest-creds x:"$TOKEN" \
              --format v2s2

            echo "Deploying to Fly.io..."
            ${pkgs.flyctl}/bin/fly deploy \
              --remote-only \
              -c "$CONFIG_FILE" \
              -i "$REGISTY" \
              -t "$TOKEN"
          '';
      };
    });
}
