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
          runtimePathOnly = true; # Only include in PATH, not closure
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
          runtimePathOnly = true;
        };
        clean = {
          exec = ''git clean -fdx'';
          description = "Clean Project";
          deps = [pkgs.git];
          runtimePathOnly = true;
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
          runtimePathOnly = true;
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
          runtimePathOnly = true;
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
          runtimePathOnly = true;
        };
        generate-docs = {
          exec = ''
            doppler run -- update -jobs 20
          '';
          deps = [
            pkgs.doppler
          ];
          extraDeps = [
            self.packages."${system}".update-slim
          ];
          description = "Update the generated go files from the md docs.";
          runtimePathOnly = true;
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
            pkgs.nix
          ];
          extraDeps = [
            scriptPackages.generate-docs
            scriptPackages.generate-css
          ];
          description = "Code Generation Steps for specific directory changes.";
          runtimePathOnly = true;
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
          runtimePathOnly = true;
        };
        generate-all = {
          exec = ''
            generate-css
            generate-docs
          '';
          deps = [];
          extraDeps = [
            scriptPackages.generate-css
            scriptPackages.generate-docs
          ];
          description = "Generate all files in parallel";
          runtimePathOnly = true;
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
            pkgs.nodePackages.prettier
          ];
          description = "Format code files";
          runtimePathOnly = true;
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
          runtimePathOnly = true;
        };
      };

      # Generate script packages with smarter dependency handling
      # Pre-define script packages to avoid dependency cycles
      scriptPackages = let
        makeScriptPackage = name: script: let
          # Determine which dependencies to include directly
          directDeps =
            if script ? runtimePathOnly && script.runtimePathOnly
            then []
            else script.deps or [];

          # Combine all deps for PATH
          allDeps = (script.deps or []) ++ (script.extraDeps or []);

          # Create PATH from all dependencies
          pathString = pkgs.lib.makeBinPath allDeps;
        in
          pkgs.writeShellApplication {
            inherit name;
            # Include PATH at the beginning of the script
            text = ''
              export PATH="${pathString}:$PATH"
              ${script.exec}
            '';
            # Only include non-runtime-path-only deps
            runtimeInputs = directDeps;
          };
      in
        builtins.listToAttrs (
          map
          (name: {
            inherit name;
            value = makeScriptPackage name scripts.${name};
          })
          (builtins.attrNames scripts)
        );

      # Build flags to strip debugging symbols and reduce binary size
      commonGoBuildFlags = [
        "-ldflags=-s -w" # Strip debug info
        "-trimpath" # Remove file paths
      ];

      # Define common attributes to avoid recursion
      goModuleCommon = {
        version = self.shortRev or "dirty";
        src = ./.;
        vendorHash = null;
      };

      # Optimized version of the conneroh package
      connerohOptimized = pkgs.buildGoModule {
        pname = "conneroh";
        inherit (goModuleCommon) version src vendorHash;

        # Only include subPackages to minimize what's built
        subPackages = ["./cmd/conneroh"];

        # Use build flags to reduce binary size
        buildFlags = commonGoBuildFlags;

        # Use upx to compress the binary if needed
        nativeBuildInputs = [
          pkgs.templ
          pkgs.tailwindcss
        ];

        preBuild = ''
          templ generate
          tailwindcss \
              --minify \
              -i ./input.css \
              -o ./cmd/conneroh/_static/dist/style.css \
              --cwd .
        '';

        # Remove UPX compression since it's causing issues
        # Instead, use build flags to reduce size

        # Don't include tests to save space
        doCheck = false;
      };

      # Optimized update binary
      updateOptimized = pkgs.buildGoModule {
        pname = "update";
        inherit (goModuleCommon) version src vendorHash;
        subPackages = ["./cmd/update"];
        buildFlags = commonGoBuildFlags;
        doCheck = false;
      };

      # Create a minimal/slim update package for use in scripts
      updateSlim = pkgs.buildGoModule {
        pname = "update-slim";
        inherit (goModuleCommon) version src vendorHash;
        subPackages = ["./cmd/update"];
        buildFlags = commonGoBuildFlags;
        doCheck = false;
        meta.mainProgram = "update"; # Ensure binary name matches
      };

      # Database files
      databaseFiles = pkgs.runCommand "database-files" {} ''
        mkdir -p $out/root
        cp ${./master.db} $out/root/master.db
      '';
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
          git status
        '';
        packages = with pkgs;
          [
            inputs.bun2nix.packages.${system}.default
            git # Essential dev tools
            go_1_24

            # Most tools added to PATH but not closure
            alejandra
            nixd
            statix
            deadnix
            air
            templ
            golangci-lint
            golines
            tailwindcss
            bun
            nodePackages.prettier
            flyctl

            # Keep only the most essential tools in closure
          ]
          # Add all script packages to PATH
          ++ builtins.attrValues scriptPackages;
      };

      # Optimized packages
      packages = let
        internal_port = 8080;
        force_https = true;
        processes = ["app"];
      in
        {
          # Optimized packages
          conneroh = connerohOptimized;
          update = updateOptimized;
          update-slim = updateSlim;

          # Optimized Docker image
          C-conneroh = pkgs.dockerTools.buildLayeredImage {
            name = "conneroh";
            tag = "latest";
            created = "now";
            config = {
              WorkingDir = "/root";
              Cmd = ["/bin/conneroh"];
              ExposedPorts = {"8080/tcp" = {};};
              Env = [
                "SSL_CERT_FILE=${pkgs.cacert}/etc/ssl/certs/ca-bundle.crt"
                "NIX_SSL_CERT_FILE=${pkgs.cacert}/etc/ssl/certs/ca-bundle.crt"
              ];
            };
            contents = [
              # Only include what's absolutely necessary
              connerohOptimized
              databaseFiles
              # Include minimal ssl certs instead of full cacert
              (pkgs.runCommand "minimal-ca-bundle" {} ''
                mkdir -p $out/etc/ssl/certs
                cp ${pkgs.cacert}/etc/ssl/certs/ca-bundle.crt $out/etc/ssl/certs/
              '')
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

            # Pre-create the toml files
            flyDevToml = settingsFormat.generate "fly.dev.toml" flyDevConfig;
            flyProdToml = settingsFormat.generate "fly.toml" flyProdConfig;

            # Capture the image path as a string variable to avoid recursion
            imageFile = "${pkgs.dockerTools.buildLayeredImage {
              name = "conneroh";
              tag = "latest";
              contents = [connerohOptimized databaseFiles];
            }}";
          in
            pkgs.writeShellApplication {
              bashOptions = [
                "errexit"
                "pipefail"
              ];
              name = "deployPackage";
              # Use PATH for dependencies instead of including them in closure
              text = ''
                # Set up PATH to include required tools
                export PATH="${pkgs.lib.makeBinPath [
                  pkgs.doppler
                  pkgs.skopeo
                  pkgs.flyctl
                  pkgs.cacert
                ]}:$PATH"

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
                  docker-archive:"${imageFile}" \
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
              # No runtime inputs - using PATH instead
              runtimeInputs = [];
            };
        }
        // pkgs.lib.genAttrs (builtins.attrNames scripts) (name: scriptPackages.${name});
    });
}
