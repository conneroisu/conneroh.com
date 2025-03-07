{
  lib,
  inputs,
  namespace,
  pkgs,
  mkShell,
  ...
}: let
  buildGoModule = pkgs.buildGoModule.override {go = pkgs.go_1_24;};
  buildWithSpecificGo = pkg: pkg.override {inherit buildGoModule;};

  # Import scripts from scripts.nix
  scripts = import ./scripts.nix {inherit lib;};

  # Convert scripts to packages
  scriptPackages =
    lib.mapAttrsToList
    (name: script: pkgs.writeShellScriptBin name script.exec)
    scripts;
in
  mkShell {
    shellHook = ''
      export REPO_ROOT=$(git rev-parse --show-toplevel)
      ${inputs.self.checks.${pkgs.system}.pre-commit.shellHook}
      export CGO_CFLAGS="-O2"

      # Print available commands
      echo "Available commands:"
      ${lib.concatStringsSep "\n" (lib.mapAttrsToList (name: script: ''echo "  ${name} - ${script.description}"'') scripts)}
    '';
    buildInputs = inputs.self.checks.${pkgs.system}.pre-commit.enabledPackages;
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
        sqlite-web

        # SQL Related
        sqlc
        sqls
        sqldiff
        inputs.sqlcquash.packages."${pkgs.system}".default
        sleek
        bc

        # Infra
        flyctl
        wireguard-tools
        openssl.dev
      ]
      # Add the generated script packages
      ++ scriptPackages
      ++ [
        pkgs."${namespace}"._copygen
      ];
  }
