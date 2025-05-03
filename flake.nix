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

  outputs = inputs @ {
    self,
    flake-utils,
    ...
  }:
    flake-utils.lib.eachDefaultSystem (system: let
      overlay = final: prev: {final.go = prev.go_1_24;};
      pkgs = import inputs.nixpkgs {
        inherit system;
        overlays = [
          overlay
        ];
        config.allowUnfree = true;
      };
      buildWithSpecificGo = pkg: pkg.override {buildGoModule = pkgs.buildGo124Module;};
      scripts = {
        dx = {
          exec = ''
            REPO_ROOT="$(git rev-parse --show-toplevel)"
            $EDITOR "$REPO_ROOT"/flake.nix
          '';
          deps = [
            pkgs.git
          ];
          description = "Edit flake.nix";
        };
        gx = {
          exec = ''
            REPO_ROOT="$(git rev-parse --show-toplevel)"
            $EDITOR "$REPO_ROOT"/go.mod
          '';
          deps = [
            pkgs.git
          ];
          description = "Edit go.mod";
        };
        clean = {
          exec = ''${pkgs.git}/bin/git clean -fdx'';
          description = "Clean Project";
        };
        reset-db = {
          exec = ''
            rm ./master.db
            rm ./master.db-shm
            rm ./master.db-wal
          '';
          description = "Reset the database";
        };
        tests = {
          exec = ''
            REPO_ROOT="$(git rev-parse --show-toplevel)"
            go test -v "$REPO_ROOT"/...
          '';
          deps = [
            pkgs.go
          ];
          description = "Run all go tests";
        };
        lint = {
          exec = ''
            REPO_ROOT="$(git rev-parse --show-toplevel)"

            golangci-lint run
            statix check "$REPO_ROOT"/flake.nix
            deadnix "$REPO_ROOT"/flake.nix
          '';
          deps = [
            pkgs.golangci-lint
            pkgs.statix
            pkgs.deadnix
          ];
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
          deps = [
            pkgs.tailwindcss
            pkgs.templ
            pkgs.go
          ];
          description = "Update the generated html and css files.";
        };
        generate-docs = {
          exec = ''
            doppler run -- update -jobs 20
          '';
          deps = [
            pkgs.doppler
            self.packages."${system}".update
          ];
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
              generate-docs
              echo "$DOCS_HASH" > ./internal/cache/docs.hash
            fi
          '';
          deps = [
            (pkgs.lib.getExe scriptPackages.generate-docs)
            (pkgs.lib.getExe scriptPackages.generate-css)
          ];
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
          deps = [
            pkgs.bun
          ];
          description = "Generate JS files";
        };
        generate-all = {
          exec = ''
            generate-css
            generate-docs
          '';
          deps = [
            (pkgs.lib.getExe scriptPackages.generate-css)
            (pkgs.lib.getExe scriptPackages.generate-docs)
          ];
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
          deps = [
            pkgs.go
            pkgs.git
            pkgs.golines
          ];
          description = "Format code files";
        };
        run = {
          exec = ''
            export DEBUG=true
            cd "$(git rev-parse --show-toplevel)" && air
          '';
          deps = [
            pkgs.air
            pkgs.git
          ];
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
      devShell = pkgs.mkShell {
        shellHook = ''
          export REPO_ROOT=$(git rev-parse --show-toplevel)
          export CGO_CFLAGS="-O2"

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
            svgcleaner
            sqlite-web
            inputs.bun2nix.packages.${system}.default
            harper

            flyctl # Infra
            openssl.dev
            skopeo
            consul

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
        version = self.shortRev or "dirty";
        src = ./.;
        vendorHash = null;
        # Create a derivation for the database file
        databaseFiles = pkgs.runCommand "database-files" {} ''
          mkdir -p $out/root
          cp ${./master.db} $out/root/master.db
        '';
      in
        rec {
          conneroh = pkgs.buildGoModule {
            inherit src vendorHash version;
            name = "conneroh.com";
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
            inherit src vendorHash version;
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
              ExposedPorts = {"8080/tcp" = {};};
              Env = [
                "SSL_CERT_FILE=${pkgs.cacert}/etc/ssl/certs/ca-bundle.crt"
                "NIX_SSL_CERT_FILE=${pkgs.cacert}/etc/ssl/certs/ca-bundle.crt"
              ];
            };
            copyToRoot = [
              conneroh
              pkgs.cacert
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
                -i "$REGISTY:latest" \
                -t "$TOKEN"
            '';
        }
        // pkgs.lib.genAttrs (builtins.attrNames scripts) (name: scriptPackages.${name});
    });
}
