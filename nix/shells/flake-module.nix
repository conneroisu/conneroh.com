{
  perSystem = {
    self',
    config,
    inputs',
    pkgs,
    system,
    ...
  }: let
    buildWithSpecificGo = pkg: pkg.override {buildGoModule = pkgs.buildGo124Module;};
    scripts = import ./shell-scripts.nix {inherit pkgs self' config;};
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
    devShells.default = pkgs.mkShellNoCC {
      shellHook = ''
        export REPO_ROOT=$(git rev-parse --show-toplevel)
        # Print available commands
        echo "Available commands:"
        ${pkgs.lib.concatStringsSep "\n" (
          pkgs.lib.mapAttrsToList (name: script: ''echo "  ${name} - ${script.description}"'') scripts
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
      packages = with pkgs;
        [
          inputs'.bun2nix.packages.default
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
    };

    packages = let
      internal_port = 8080;
      force_https = true;
      processes = ["app"];
      version = self'.shortRev or "dirty";
      src = ./../../.;
      # Create a derivation for the database file
      databaseFiles = pkgs.runCommand "database-files" {} ''
        mkdir -p $out/root
        cp ${./../../master.db} $out/root/master.db
      '';

      preBuild = ''
        ${pkgs.templ}/bin/templ generate
        ${pkgs.tailwindcss}/bin/tailwindcss \
            --minify \
            -i ./input.css \
            -o ./cmd/conneroh/_static/dist/style.css \
            --cwd .
      '';
    in
      {
        conneroh = pkgs.buildGoModule {
          inherit src version preBuild;
          vendorHash = "sha256-kOGauV5dMTcHvSR7uWvY1dcKR4WqlWccDfnXtycsRVI=";
          name = "conneroh.com";
          goSum = ./../../go.sum;
          subPackages = ["."];
        };
        update = pkgs.buildGoModule {
          inherit src version preBuild;
          vendorHash = "sha256-kOGauV5dMTcHvSR7uWvY1dcKR4WqlWccDfnXtycsRVI=";
          goSum = ./../../go.sum;
          name = "update";
          subPackages = ["./cmd/update"];
          doCheck = false;
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
            self'.packages.conneroh
            databaseFiles
          ];
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
          pkgs.writeShellApplication {
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
                docker-archive:"${self'.packages.C-conneroh}" \
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
  };
}
