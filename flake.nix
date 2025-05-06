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
        interpolate = {
          exec = ''
            FILE="$1"
            START_MARKER="$2"
            END_MARKER="$3"
            NEW_CONTENT="$4"

            # Create a temporary file with mktemp (safer than fixed filenames)
            TEMP_FILE=$(mktemp)

            # Write the new content to the temporary file
            echo "$NEW_CONTENT" > "$TEMP_FILE"

            # Perform the replacement
            sed -i -e "/$START_MARKER/,/$END_MARKER/{
                /$START_MARKER/{
                    p
                    r $TEMP_FILE
                }
                /$END_MARKER/p
                d
            }" "$FILE"

            # Clean up
            rm "$TEMP_FILE"
          '';
          deps = with pkgs; [templ];
          description = "Interpolate templates; Usage: interpolate input_file start_marker end_marker replacement_text";
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
        generate-docs = {
          exec = ''
            REPO_ROOT="$(git rev-parse --show-toplevel)"

            # Create a temporary file with the content to insert
            TEMP_CONTENT=$(mktemp)

            # Generate all content with a single redirection
            {
              ${builtins.concatStringsSep "\n" (
              pkgs.lib.mapAttrsToList (
                name: script: ''echo "  ${name} - ${script.description}"''
              )
              scripts
            )}
            } > "$TEMP_CONTENT"

            # Use the interpolate command with the content from the file
            interpolate "$REPO_ROOT"/README.md "<!-- BEGIN_MARKER -->" "<!-- END_MARKER -->" "$(cat "$TEMP_CONTENT")"

            # Clean up
            rm "$TEMP_CONTENT"
          '';
          deps = with pkgs; [doppler coreutils self.packages."${system}".interpolate];
          description = "Update the generated documentation files.";
        };
        generate-db = {
          exec = ''
            doppler run -- update
          '';
          deps = with pkgs; [doppler self.packages."${system}".update];
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
        run = {
          exec = ''
            export DEBUG=true
            cd "$(git rev-parse --show-toplevel)" && air
          '';
          deps = with pkgs; [air git];
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
      devShells.default = pkgs.mkShell {
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
          git status
        '';
        packages = with pkgs;
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
        vendorHash = "sha256-BUI6XA3RnQWKrNohX1iC3UxXpc+9WcHxrnX+bxgEpTU=";
        # Create a derivation for the database file
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
      in
        {
          conneroh = pkgs.buildGoModule {
            inherit src vendorHash version preBuild;
            name = "conneroh.com";
            nativeBuildInputs = [
              pkgs.templ
              pkgs.tailwindcss
            ];
            subPackages = ["."];
          };
          update = pkgs.buildGoModule {
            inherit src vendorHash version preBuild;
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
              self.packages.${system}.conneroh
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
              bashOptions = [
                "errexit"
                "pipefail"
              ];
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

                REGISTY="registry.fly.io/$FLY_NAME"
                echo "Copying image to Fly.io... to $REGISTY"

                skopeo copy \
                  --insecure-policy \
                  docker-archive:"${self.packages.${system}.C-conneroh}" \
                  "docker://$REGISTY:latest" \
                  --dest-creds x:"$TOKEN" \
                  --format v2s2

                echo "Deploying to Fly.io..."
                fly deploy \
                  --remote-only \
                  -c "$CONFIG_FILE" \
                  -i "$REGISTY:latest" \
                  -t "$TOKEN"
              '';
            };
        }
        // pkgs.lib.genAttrs (builtins.attrNames scripts) (name: scriptPackages.${name});
    });
}
