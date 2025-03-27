{
  description = "Description for the project";

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

    mk-shell-bin.url = "github:rrbutani/nix-mk-shell-bin";

    sqlcquash = {
      url = "github:conneroisu/sqlcquash/main";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.flake-utils.follows = "flake-utils";
    };

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

  outputs = inputs @ {flake-utils, ...}:
    flake-utils.lib.eachSystem [
      "x86_64-linux"
      "i686-linux"
      "x86_64-darwin"
      "aarch64-linux"
      "aarch64-darwin"
    ] (system: let
      overlays = [(final: prev: {go = prev.go_1_24;})];
      pkgs = import inputs.nixpkgs {inherit system overlays;};
      buildGoModule = pkgs.buildGoModule.override {go = pkgs.go_1_24;};
      buildWithSpecificGo = pkg: pkg.override {inherit buildGoModule;};

      scripts = {
        dx = {
          exec = ''$EDITOR $REPO_ROOT/flake.nix'';
          description = "Edit the flake.nix";
        };
        clean = {
          exec = ''
            git clean -fdx
          '';
          description = "Clean Project";
        };
        tests = {
          exec = ''go test -v -short ./...'';
          description = "Run go tests with short flag";
        };
        unit-tests = {
          exec = ''go test -v ./...'';
          description = "Run all go tests";
        };
        lint = {
          exec = ''golangci-lint run'';
          description = "Run golangci-lint";
        };
        build = {
          exec = ''nix build --accept-flake-config .#packages.x86_64-linux.conneroh'';
          description = "Build the package";
        };
        update = {
          exec = ''go run $REPO_ROOT/cmd/update'';
          description = "Update the database.";
        };
        restart = {
          exec = ''rm -f $REPO_ROOT/master.db && go run $REPO_ROOT/cmd/update'';
          description = "Execute restart command with doppler";
        };
        generate-reload = {
          exec = ''
            templ generate $REPO_ROOT &
            wait
          '';
          description = "Generate templ files and wait for completion";
        };
        generate-js = {
          exec = ''
            export REPO_ROOT=$(git rev-parse --show-toplevel) # needed
            bun build \
                $REPO_ROOT/index.js \
                --minify \
                --minify-syntax \
                --minify-whitespace  \
                --minify-identifiers \
                --outdir $REPO_ROOT/cmd/conneroh/_static/dist/ &
          '';
          description = "Generate js files";
        };

        generate-all = {
          exec = ''
            go generate $REPO_ROOT/... &

            templ generate $REPO_ROOT &

            generate-js &
            tailwindcss \
                --minify \
                -i $REPO_ROOT/input.css \
                -o $REPO_ROOT/cmd/conneroh/_static/dist/style.css \
                --cwd $REPO_ROOT &

            wait
          '';
          description = "Generate all files in parallel";
        };
        format = {
          exec = ''
            export REPO_ROOT=$(git rev-parse --show-toplevel) # needed

            go fmt $REPO_ROOT/...

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
    in {
      devShells.default = pkgs.mkShell {
        shellHook = ''
          export REPO_ROOT=$(git rev-parse --show-toplevel)
          export CGO_CFLAGS="-O2"

          # Print available commands
          echo "Available commands:"
          ${pkgs.lib.concatStringsSep "\n" (pkgs.lib.mapAttrsToList (name: script: ''echo "  ${name} - ${script.description}"'') scripts)}
        '';
        packages = with pkgs;
          [
            # Nix
            alejandra
            nixd

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

            # Web
            tailwindcss
            tailwindcss-language-server
            bun
            nodePackages.typescript-language-server
            nodePackages.prettier
            sqlite-web
            nodePackages.svgo

            # SQL Related
            sqlc
            sqls
            sqldiff
            inputs.sqlcquash.packages."${pkgs.system}".default
            sleek
            bc

            # C/C++
            clang-tools

            # Infra
            flyctl
            wireguard-tools
            openssl.dev
          ]
          # Add the generated script packages
          ++ scriptPackages;
      };

      packages = let
        app-name = "conneroh.com";
      in rec {
        conneroh = buildGoModule {
          # pname = app-name;
          name = app-name;
          version = "0.0.1";
          src = ./.;
          subPackages = ["."];
          nativeBuildInputs = [];
          vendorHash = "sha256-KeUHn4w8Xc0He/mg6XJoK+0276WiTw5phIquwY7Usaw=";
          preBuild = ''
            ${pkgs.templ}/bin/templ generate
          '';
        };
        C-conneroh = pkgs.dockerTools.buildLayeredImage {
          # pname = app-name;
          name = app-name;
          tag = "latest";
          contents = [
            conneroh
            pkgs.cacert
          ];
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
          extraCommands = ''
            echo "$(git rev-parse HEAD)" > REVISION
          '';
        };
        deployPackage = pkgs.writeShellScriptBin "deploy" ''
          set -e

          # Check if FLY_AUTH_TOKEN is set
          if [ -z ''${FLY_AUTH_TOKEN+x} ]; then
            echo "Error: FLY_AUTH_TOKEN environment variable is not set"
            echo "Please set it with: export FLY_AUTH_TOKEN=your-token"
            exit 1
          fi

          echo "Copying image to Fly.io registry..."
          ${pkgs.skopeo}/bin/skopeo --insecure-policy copy \
            docker-archive:"$IMAGE_PATH" \
            docker://registry.fly.io/$APP_NAME:latest \
            --dest-creds x:"$FLY_AUTH_TOKEN" \
            --format v2s2

          echo "Triggering deployment..."
          ${pkgs.curl}/bin/curl -X POST \
            -H "Authorization: Bearer $FLY_AUTH_TOKEN" \
            -H "Content-Type: application/json" \
            "https://api.fly.io/api/v1/apps/$APP_NAME/deploys" \
            -d "{\"image\":\"registry.fly.io/$APP_NAME:latest\"}"

          echo "Deployment initiated!"
        '';
      };
    });
}
