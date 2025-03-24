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
  scripts = import ./scripts.nix {inherit lib pkgs;};

  # Convert scripts to packages
  scriptPackages =
    lib.mapAttrsToList
    (name: script: pkgs.writeShellScriptBin name script.exec)
    scripts;
in
  mkShell {
    shellHook = ''
      export REPO_ROOT=$(git rev-parse --show-toplevel)
      export CGO_CFLAGS="-O2"

      # Print available commands
      echo "Available commands:"
      ${lib.concatStringsSep "\n" (lib.mapAttrsToList (name: script: ''echo "  ${name} - ${script.description}"'') scripts)}
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
        llama-cpp
      ]
      # Add the generated script packages
      ++ scriptPackages;
  }
