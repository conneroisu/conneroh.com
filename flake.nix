{
  description = "Personal Website for Conner Ohnesorge";

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
    flake-utils.lib.eachSystem ["x86_64-linux" "aarch64-linux" "aarch64-darwin"] (
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
        flake-scripts = import ./flake-scripts.nix {inherit pkgs self system;};
        scriptPackages =
          pkgs.lib.mapAttrs
          (
            name: script:
              pkgs.writeShellApplication {
                inherit name;
                text = script.exec;
                runtimeInputs = script.deps or [];
              }
          )
          flake-scripts.scripts;
      in {
        devShells = let
          shell-shellHook = ''
            export REPO_ROOT=$(git rev-parse --show-toplevel)
            echo "Available commands:"
            ${pkgs.lib.concatStringsSep "\n" (
              pkgs.lib.mapAttrsToList (name: script: ''echo "  ${name} - ${script.description}"'') flake-scripts.scripts
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
                # Container Deps
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
            path = pkgs.lib.cleanSourceWith {
              src = ./.;
              filter = path: type: let
                baseName = baseNameOf path;
              in
                !(
                  (type == "directory" && baseName == ".direnv")
                  || (type == "symlink" && baseName == ".git")
                  || (type == "directory" && baseName == "data")
                );
            };
            name = "source";
          };

          databaseFiles = pkgs.runCommand "database-files" {} ''
            mkdir -p $out/root
            # Joint Shm and Wal
            ${pkgs.sqlite}/bin/sqlite3 ${./master.db} "PRAGMA wal_checkpoint(FULL);"
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
          // pkgs.lib.genAttrs (builtins.attrNames flake-scripts.scripts) (name: scriptPackages.${name});
      }
    );
}
