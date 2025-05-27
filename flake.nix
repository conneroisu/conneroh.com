{
  description = "Personal Website for Conner Ohnesorge";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    systems.url = "github:nix-systems/default";
    bun2nix.url = "github:baileyluTCD/bun2nix";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = inputs @ {
    self,
    flake-utils,
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
              alejandra # Nix
              nixd
              nil
              statix
              deadnix
              inputs.bun2nix.packages.${system}.default

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
              git-bug

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
        };

        apps = {
          pr-preview = {
            type = "app";
            meta.description = "A preview server for pull requests";
            program = "${self.packages.${system}.pr-preview}/bin/pr-preview";
          };
          deployPackage = {
            type = "app";
            meta.description = "Deploys the conneroh.com Docker image to Fly.io";
            program = "${self.packages.${system}.deployPackage}/bin/deployPackage";
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

          flyProdToml = settingsFormat.generate "fly.toml" flyProdConfig;
        in
          {
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

                [ -z "$FLY_AUTH_TOKEN" ] && FLY_AUTH_TOKEN="$(doppler secrets get --plain FLY_AUTH_TOKEN)"
                TOKEN="$FLY_AUTH_TOKEN"
                export FLY_NAME="conneroh-com"
                export CONFIG_FILE=${flyProdToml}

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

            # PR Preview deployment script
            pr-preview = pkgs.writeShellScriptBin "pr-preview" ''
              set -euo pipefail

              # Add required tools to PATH
              export PATH="${
                pkgs.lib.makeBinPath (with pkgs; [flyctl skopeo jq git gnused coreutils])
              }:$PATH"

              # Script configuration
              readonly APP_PREFIX="pr"
              readonly FLY_ORG="''${FLY_ORG:-personal}"
              readonly FLY_REGION="''${FLY_REGION:-ord}"

              # Functions
              generate_app_name() {
                  local pr_number="$1"
                  echo "''${APP_PREFIX}-''${pr_number}-conneroh-com" | tr '[:upper:]' '[:lower:]'
              }

              destroy_pr_app() {
                  local pr_number="$1"

                  local app_name
                  app_name=$(generate_app_name "$pr_number")

                  echo "Destroying app: ''${app_name}"

                  if flyctl apps list --json | jq -e ".[] | select(.Name == \"''${app_name}\")" > /dev/null; then
                      flyctl apps destroy "''${app_name}" --yes
                      echo "App ''${app_name} destroyed successfully"
                  else
                      echo "App ''${app_name} not found, nothing to destroy"
                  fi
              }

              deploy_pr_app() {
                  local pr_number="$1"
                  shift

                  local app_name
                  app_name=$(generate_app_name "$pr_number")

                  echo "Deploying PR #''${pr_number} to app: ''${app_name}"

                  # Check if app exists
                  if ! flyctl apps list --json | jq -e ".[] | select(.Name == \"''${app_name}\")" > /dev/null; then
                      echo "Creating new app: ''${app_name}"
                      flyctl apps create "''${app_name}" --org "''${FLY_ORG}"
                  fi

                  # Create fly.toml for PR preview
                  cat > fly.pr.toml <<EOF
              app = "''${app_name}"
              primary_region = "''${FLY_REGION}"

              [http_service]
                internal_port = ${toString internal_port}
                force_https = ${
                if force_https
                then "true"
                else "false"
              }
                auto_stop_machines = "stop"
                auto_start_machines = true
                min_machines_running = 0
                processes = ["app"]

              [[vm]]
                memory = "512M"
                cpu_kind = "shared"
                cpus = 1
              EOF

                  # Copy image to Fly.io registry
                  local registry="registry.fly.io/''${app_name}"
                  echo "Copying image to ''${registry}..."

                  skopeo copy \
                    --insecure-policy \
                    docker-archive:"${self.packages."${system}".C-conneroh}" \
                    "docker://''${registry}:latest" \
                    --dest-creds x:"''${MASTER_FLY_AUTH_TOKEN}"

                  # Deploy
                  flyctl deploy \
                    --app "''${app_name}" \
                    --config fly.pr.toml \
                    --image "''${registry}:latest" \
                    --remote-only \
                    "$@"

                  # Output deployment information
                  echo "Deployment complete!"
                  echo "URL: https://''${app_name}.fly.dev"

                  # Get deployment details
                  flyctl status --app "''${app_name}" --json | jq '{
                      app: .Name,
                      url: "https://\(.Name).fly.dev",
                      version: .DeploymentStatus.Version,
                      status: .DeploymentStatus.Status
                  }'
              }

              # Main command handling
              case "''${1:-}" in
                  deploy)
                      shift
                      deploy_pr_app "$@"
                      ;;
                  destroy)
                      shift
                      destroy_pr_app "$@"
                      ;;
                  *)
                      echo "Usage: pr-preview {deploy|destroy} <pr_number> [additional args]"
                      exit 1
                      ;;
              esac
            '';
          }
          // pkgs.lib.genAttrs (builtins.attrNames flake-scripts.scripts) (
            name: scriptPackages.${name}
          );
      }
    );
}
