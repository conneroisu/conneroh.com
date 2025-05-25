{
  description = "Personal Website for Conner Ohnesorge";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    systems.url = "github:nix-systems/default";
    bun2nix.url = "github:baileyluTCD/bun2nix";
    flake-utils.url = "github:numtide/flake-utils";
    msb.url = "github:rrbutani/nix-mk-shell-bin";
  };

  outputs = inputs @ {
    self,
    flake-utils,
    msb,
    ...
  }:
    flake-utils.lib.eachSystem ["x86_64-linux" "aarch64-linux" "aarch64-darwin"] (
      system: let
        pkgs = import inputs.nixpkgs {
          inherit system;
          overlays = [
            (final: prev: {final.go = prev.go_1_24;})
          ];
          config.allowUnfree = true;
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
          shellHook = ''
            export REPO_ROOT=$(git rev-parse --show-toplevel)
            echo "Available commands:"
            ${pkgs.lib.concatStringsSep "\n" (
              pkgs.lib.mapAttrsToList (name: script: ''echo "  ${name} - ${script.description}"'') flake-scripts.scripts
            )}
          '';

          env = {
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
              harper
              htmx-lsp
              vscode-langservers-extracted
              sqlite

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
            inherit shellHook env;
            packages = shell-packages;
          };
          devcontainer = pkgs.mkShell {
            inherit shellHook env;
            packages =
              shell-packages
              ++ (with pkgs; [
                # Container Deps
                nix
                coreutils-full
                toybox
                curl
                getent # Required by devcontainers for user listing
                bashInteractive # Full bash shell (not just sh)
                shadow # User management utilities
                sudo
                wget
                docker
                git
                gnugrep
                gnused
                jq
                skopeo
                util-linux
                gh
                vscode
                code-server
                gnugrep
                gnused
              ]);
          };
        };

        packages = let
          internal_port = 8080;
          force_https = true;
          processes = ["app"];
          version = self.shortRev or "dirty";
          src = pkgs.lib.cleanSourceWith {
            src = ./.;
            filter = path: type: let
              baseName = baseNameOf path;
              path' = toString path;
              isExcluded =
                (baseName == ".direnv")
                || (baseName == ".git")
                || (baseName == "node_modules")
                || (baseName == "data" && type == "directory")
                || (builtins.match ".*/internal/data(/.*|$)" path' != null)
                || (baseName == "result")
                || (pkgs.lib.hasSuffix ".swp" baseName)
                || (pkgs.lib.hasSuffix "~" baseName);
            in
              !isExcluded;
            name = "source";
          };

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
          alpine = pkgs.dockerTools.pullImage {
            imageName = "alpine";
            imageDigest = "sha256:a8560b36e8b8210634f77d9f7f9efd7ffa463e380b75e2e74aff4511df3ef88c";
            arch = "amd64";
            sha256 = "sha256-Eb4oYIKZOj6lg8ej+/4sFFCvvJtrzwjKRjtBQG8CHJQ=";
          };
          tag = "v8";
        in
          {
            devShellBin = msb.lib.mkShellBin {
              drv = self.devShells.${system}.devcontainer;
              nixpkgs = pkgs;
              bashPrompt = "[conneroh]$ ";
            };
            devContainer = pkgs.dockerTools.buildImage {
              name = "devContainer";
              fromImage = alpine;
              tag = "latest";
              runAsRoot = ''
                #!${pkgs.runtimeShell}
                ${pkgs.dockerTools.shadowSetup}
                # Create vscode group and user with UID 1000
                groupadd -r -g 1000 vscode
                useradd -r -g vscode -u 1000 vscode -m -s ${pkgs.bashInteractive}/bin/bash
                # Create home directory with proper permissions
                mkdir -p /home/vscode
                chown -R vscode:vscode /home/vscode
                # Add sudo access
                mkdir -p /etc/sudoers.d
                echo "vscode ALL=(ALL) NOPASSWD:ALL" > /etc/sudoers.d/vscode
                chmod 0440 /etc/sudoers.d/vscode
              '';
              config = {
                entrypoint = [
                  "${self.packages.${system}.devShellBin}/bin/nix-shell-env-shell"
                ];
                User = "vscode";
                WorkingDir = "/home/vscode";
                Env = [
                  "USER=vscode"
                  "HOME=/home/vscode"
                ];
              };
            };
            conneroh = pkgs.buildGoModule {
              inherit src version preBuild;
              vendorHash = "sha256-DYqIBhMpuNc62m9fCU7T6Sl17tmpTztD70qG1OGUEN8=";
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
                (pkgs.runCommand "database-files" {} ''
                  mkdir -p $out/root
                  cp ${./master.db} $out/root/master.db
                '')
              ];
            };
            deployPackage = pkgs.writeShellApplication {
              name = "deployPackage";
              runtimeInputs = with pkgs; [doppler skopeo flyctl cacert];
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
            deployDev = pkgs.writeShellApplication {
              name = "deployDev";
              runtimeInputs = with pkgs; [doppler skopeo flyctl cacert];
              bashOptions = ["errexit" "pipefail"];
              text = ''
                set -e
                [ -z "$GHCR_TOKEN" ] && GHCR_TOKEN="$(doppler secrets get --plain GHCR_TOKEN)"
                TOKEN="$GHCR_TOKEN"

                REGISTRY="ghcr.io/conneroisu/conneroh.com"
                echo "Copying image to $REGISTRY"
                skopeo copy \
                  --insecure-policy \
                  docker-archive:"${self.packages."${system}".devContainer}" \
                  "docker://$REGISTRY:${tag}" \
                  --dest-creds x:"$TOKEN" \
                  --format v2s2
              '';
            };
          }
          // pkgs.lib.genAttrs (builtins.attrNames flake-scripts.scripts) (
            name: scriptPackages.${name}
          );
      }
    );
}
