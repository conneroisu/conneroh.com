{
  description = "Personal Website for Conner Ohnesorge";

  # TODO: Might be adventageous to introduce flake-utils to reduce copy-pasta of nix paths
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    systems.url = "github:nix-systems/default";
    bun2nix.url = "github:baileyluTCD/bun2nix";
    flake-utils.url = "github:numtide/flake-utils";
    nix2container.url = "github:nlewo/nix2container";
    nix2container.inputs = {
      nixpkgs.follows = "nixpkgs";
    };
    rust-overlay.url = "github:oxalica/rust-overlay";
    rust-overlay.inputs.nixpkgs.follows = "nixpkgs";
  };

  outputs = inputs @ {
    flake-utils,
    self,
    ...
  }:
    flake-utils.lib.eachSystem ["x86_64-linux" "aarch64-linux"] (
      system: let
        pkgs = import inputs.nixpkgs {
          inherit system;
          config.allowUnfree = true;
          overlays = [
            (final: prev: {
              go = prev.go_1_24;
            })
            inputs.rust-overlay.overlays.default
          ];
        };

        buildWithSpecificGo = pkg: pkg.override {buildGoModule = pkgs.buildGo124Module;};
        scripts = {
          dx = {
            exec = ''$EDITOR "$(git rev-parse --show-toplevel)"/flake.nix'';
            deps = [pkgs.git];
            description = "Edit flake.nix";
          };
          gx = {
            exec = ''$EDITOR "$(git rev-parse --show-toplevel)"/go.mod'';
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
              doppler run -- go run ./cmd/update
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
            deps = with self.packages."${system}"; [generate-db generate-css];
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
      in {
        devShells = let
          shell-shellHook = ''
            export REPO_ROOT=$(git rev-parse --show-toplevel)
            # Print available commands
            echo "Available commands:"
            ${pkgs.lib.concatStringsSep "\n" (
              pkgs.lib.mapAttrsToList (name: script: ''echo "  ${name} - ${script.description}"'') scripts
            )}
          '';

          shell-env = {
            PLAYWRIGHT_BROWSERS_PATH = "${pkgs.playwright-driver.browsers}";
            PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD = "1";
            PLAYWRIGHT_NODEJS_PATH = "${pkgs.nodejs_20}/bin/node";

            # Browser executable paths
            PLAYWRIGHT_CHROMIUM_EXECUTABLE_PATH = "${ #
              "${pkgs.playwright-driver.browsers}/chromium-1155"
            }";
          };
          shell-packages = with pkgs;
            [
              inputs.bun2nix.packages.${system}.default
              alejandra # Nix
              nixd
              nil
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
              yaml-language-server
              nodePackages.typescript-language-server
              nodePackages.prettier
              svgcleaner
              sqlite-web
              harper
              htmx-lsp
              vscode-langservers-extracted
              sqlite

              rust-bin.stable.latest.default # Rust
              rust-analyzer
              pkg-config

              flyctl # Infra
              openssl.dev
              skopeo

              (
                pkgs.buildGoModule rec {
                  pname = "copygen";
                  version = "0.4.1";

                  src = pkgs.fetchFromGitHub {
                    owner = "switchupcb";
                    repo = "copygen";
                    rev = "v${version}";
                    sha256 = "sha256-gdoUvTla+fRoYayUeuRha8Dkix9ACxlt0tkac0CRqwA=";
                  };

                  vendorHash = "sha256-dOIGGZWtr8F82YJRXibdw3MvohLFBQxD+Y4OkZIJc2s=";
                  subPackages = ["."];
                  proxyVendor = true;

                  ldflags = [
                    "-s"
                    "-w"
                    "-X main.version=${version}"
                  ];

                  meta = with lib; {
                    description = "Copygen";
                    homepage = "https://github.com/switchupcb/copygen";
                    license = licenses.mit;
                    mainProgram = "copygen";
                  };
                }
              )
            ]
            ++ builtins.attrValues scriptPackages;
        in {
          default = pkgs.mkShell {
            shellHook = shell-shellHook;
            env = shell-env;
            packages = shell-packages;
          };
          devcontainer = pkgs.mkShell {
            env = shell-env;
            shellHook = shell-shellHook;
            packages =
              shell-packages
              ++ (with pkgs; [
                # Contaier Deps
                coreutils-full
                toybox
                curl
                wget
                docker
                git
                gnugrep
                gnused
                jq
                nix
                skopeo
                util-linux
                gh
              ]);
          };
        };

        packages = let
          internal_port = 8080;
          force_https = true;
          processes = ["app"];
          version = self.shortRev or "dirty";
          src = builtins.path {
            path = ./.;
            name = "source";
            filter = path: type: !builtins.match "internal/docs" path;
          };

          databaseFiles = pkgs.runCommand "database-files" {} ''
            mkdir -p $out/root
            cp ${./master.db} $out/root/master.db
          '';

          preBuild = ''
            ${pkgs.templ}/bin/templ generate
            ${pkgs.tailwindcss}/bin/tailwindcss \
                --minify \
                -i ./input.css \
                -o ./cmd/conneroh/_static/dist/style.css \
                --cwd .
          '';

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
          {
            conneroh = pkgs.buildGoModule {
              inherit src version preBuild;
              vendorHash = "sha256-bPcOM7B+17SBqrEfAdJdUDEqkYWzlvys8YD1gh1mbX8=";
              name = "conneroh.com";
              goSum = ./go.sum;
              subPackages = ["."];
            };
            C-conneroh = pkgs.dockerTools.buildImage {
              created = "now";
              tag = "latest";
              name = "conneroh";
              config = {
                WorkingDir = "/root";
                Cmd = ["/bin/conneroh.com"];
                ExposedPorts = {"8080/tcp" = {};};
                Env = [
                  "SSL_CERT_FILE=${pkgs.cacert}/etc/ssl/certs/ca-bundle.crt"
                  "NIX_SSL_CERT_FILE=${pkgs.cacert}/etc/ssl/certs/ca-bundle.crt"
                ];
              };
              copyToRoot = [
                self.packages."${system}".conneroh
                databaseFiles
              ];
            };
            deployPackage = pkgs.writeShellApplication {
              name = "deployPackage";
              runtimeInputs = with pkgs; [
                doppler
                skopeo
                flyctl
                cacert
              ];
              bashOptions = ["errexit" "pipefail"];
              text = ''
                set -e
                arg=$1
                TOKEN=""
                FLY_NAME=""
                CONFIG_FILE=""

                [ -z "$arg" ] && { echo "No argument provided. Please provide a target environment. (dev or prod)"; exit 1; }

                if [ "$arg" = "dev" ]; then
                  [ -z "$FLY_DEV_AUTH_TOKEN" ] && FLY_DEV_AUTH_TOKEN="$(doppler secrets get --plain FLY_DEV_AUTH_TOKEN)"
                  TOKEN="$FLY_DEV_AUTH_TOKEN"
                  export FLY_NAME="conneroh-com-dev"
                  export CONFIG_FILE=${flyDevToml}
                else
                  [ -z "$FLY_AUTH_TOKEN" ] && FLY_AUTH_TOKEN="$(doppler secrets get --plain FLY_AUTH_TOKEN)"
                  TOKEN="$FLY_AUTH_TOKEN"
                  export FLY_NAME="conneroh-com"
                  export CONFIG_FILE=${flyProdToml}
                fi

                REGISTRY="registry.fly.io/$FLY_NAME"
                echo "Copying image to Fly.io... to $REGISTRY"

                skopeo copy \
                  --insecure-policy \
                  docker-archive:"${self.packages."${system}".C-conneroh}" \
                  "docker://$REGISTRY:latest" \
                  --dest-creds x:"$TOKEN" \
                  --format v2s2

                echo "Deploying to Fly.io..."
                fly deploy \
                  --remote-only \
                  -c "$CONFIG_FILE" \
                  -i "$REGISTRY:latest" \
                  -t "$TOKEN"
              '';
            };
          }
          // pkgs.lib.genAttrs (builtins.attrNames scripts) (name: scriptPackages.${name});
      }
    );
}
