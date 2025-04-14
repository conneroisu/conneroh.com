{
  description = "Programatically generate static structs in Go.";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";

    flake-parts.url = "github:hercules-ci/flake-parts";
    flake-parts.inputs.nixpkgs-lib.follows = "nixpkgs";

    flake-utils.url = "github:numtide/flake-utils";
    flake-utils.inputs.systems.follows = "systems";

    systems.url = "github:nix-systems/default";
  };

  nixConfig = {
    extra-substituters = ''
      https://cache.nixos.org
      https://nix-community.cachix.org
      https://devenv.cachix.org
    '';
    extra-trusted-public-keys = ''
      cache.nixos.org-1:6NCHdD59X431o0gWypbMrAURkbJ16ZPMQFGspcDShjY=
      nix-community.cachix.org-1:mB9FSh9qf2dCimDSUo8Zy7bkq5CX+/rkCWyvRCYg3Fs=
      devenv.cachix.org-1:w1cLUi8dv3hnoSPGAuibQv+f9TZLr6cv/Hm9XgU50cw=
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
          description = "Edit flake.nix";
        };
        tests = {
          exec = ''go test -v -short ./...'';
          description = "Run short go tests";
        };
        unit-tests = {
          exec = ''go test -v ./...'';
          description = "Run all go tests";
        };
        coverage-tests = {
          exec = ''go test -coverprofile=coverage.out ./...'';
          description = "Run all go tests with coverage";
        };
        lint = {
          exec = ''golangci-lint run'';
          description = "Run golangci-lint";
        };
        generate-all = {
          exec = ''
            go generate -v "$REPO_ROOT/..."
            format
            echo "All tasks completed!"
          '';
          description = "Generate all code artifacts";
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
              | xargs prettierd --write

            golines -l -w --max-len=79 --shorten-comments  --ignored-dirs=.devenv .

          '';
          description = "Format code files across multiple languages";
        };
      };

      # Convert scripts to packages
      scriptPackages =
        pkgs.lib.mapAttrsToList
        (name: script: pkgs.writeShellScriptBin name script.exec)
        scripts;
    in {
      packages = {
        doc = pkgs.stdenv.mkDerivation {
          pname = "genstruct-docs";
          version = "0.1";
          src = ./.;
          nativeBuildInputs = with pkgs; [
            nixdoc
            mdbook
            mdbook-open-on-gh
            mdbook-cmdrun
            git
          ];
          dontConfigure = true;
          dontFixup = true;
          env.RUST_BACKTRACE = 1;
          buildPhase = ''
            runHook preBuild
            cd doc  # Navigate to the doc directory during build
            mkdir -p .git  # Create .git directory
            mdbook build
            runHook postBuild
          '';
          installPhase = ''
            runHook preInstall
            mv book $out
            runHook postInstall
          '';
        };
      };

      devShells.default = pkgs.mkShell {
        shellHook = ''
          export REPO_ROOT=$(git rev-parse --show-toplevel)
          export CGO_CFLAGS="-O2"
        '';
        packages = with pkgs;
          [
            # Nix
            alejandra
            nixd

            # Go Tools
            go_1_24
            air
            pprof
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

            # Formatters
            prettierd
          ]
          ++ scriptPackages;
      };
    });
}
