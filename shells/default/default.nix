{
  # Snowfall Lib provides a customized `lib` instance with access to your flake's library
  # as well as the libraries available from your flake's inputs.
  lib,
  # You also have access to your flake's inputs.
  inputs,
  # The namespace used for your flake, defaulting to "internal" if not set.
  namespace,
  # All other arguments come from NixPkgs. You can use `pkgs` to pull shells or helpers
  # programmatically or you may add the named attributes as arguments here.
  pkgs,
  mkShell,
  ...
}: let
  buildGoModule = pkgs.buildGoModule.override {go = pkgs.go_1_24;};
  buildWithSpecificGo = pkg: pkg.override {inherit buildGoModule;};
in
  mkShell {
    shellHook = ''
      export REPO_ROOT=$(git rev-parse --show-toplevel)
      ${inputs.self.checks.${pkgs.system}.pre-commit.shellHook}
    '';
    buildInputs = inputs.self.checks.${pkgs.system}.pre-commit.enabledPackages;
    packages = with pkgs; [
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

      # SQL Related
      sqlc
      sqls
      sqldiff
      inputs.sqlcquash.packages."${pkgs.system}".default
      sleek
      nodePackages.typescript-language-server

      # Infra
      flyctl
      wireguard-tools
      openssl.dev

      (pkgs.writeShellScriptBin "dx" ''$EDITOR $REPO_ROOT/flake.nix'')
      (pkgs.writeShellScriptBin "tests" ''go test -v -short ./...'')
      (pkgs.writeShellScriptBin "unit-tests" ''go test -v ./...'')
      (pkgs.writeShellScriptBin "lint" ''golangci-lint run'')
      (pkgs.writeShellScriptBin "build" ''nix build .#packages.x86_64-linux.conneroh'')
      (pkgs.writeShellScriptBin "generate-reload" ''
        bun build \
         $REPO_ROOT/index.js \
         --minify \
         --minify-syntax \
         --minify-whitespace  \
         --minify-identifiers \
         --outdir $REPO_ROOT/cmd/conneroh/_static/dist/
        tailwindcss \
         --minify \
         -i $REPO_ROOT/input.css \
         -o $REPO_ROOT/cmd/conneroh/_static/dist/style.css \
         --cwd $REPO_ROOT
      '')
      (pkgs.writeShellScriptBin "generate-all" ''
        go generate $REPO_ROOT/... &
        templ generate $REPO_ROOT &
        bun build \
            $REPO_ROOT/index.js \
            --minify \
            --minify-syntax \
            --minify-whitespace  \
            --minify-identifiers \
            --outdir $REPO_ROOT/cmd/conneroh/_static/dist/ &
        tailwindcss \
            --minify \
            -i $REPO_ROOT/input.css \
            -o $REPO_ROOT/cmd/conneroh/_static/dist/style.css \
            --cwd $REPO_ROOT &

          wait
      '')
      (pkgs.writeShellScriptBin "format" ''
        export REPO_ROOT=$(git rev-parse --show-toplevel) # needed
        go fmt $REPO_ROOT/...

        git ls-files \
            --others \
            --exclude-standard \
            --cached \
            -- '*.cc' '*.h' '*.proto' \
            | xargs clang-format -i

        git ls-files \
          --others \
          --exclude-standard \
          --cached \
          -- '*.js' '*.ts' '*.css' '*.md' '*.json' \
          | xargs prettier --write

        golines -l -w --max-len=80 --shorten-comments  --ignored-dirs=.devenv .

      '')

      (pkgs.writeShellScriptBin "run" ''air'')
    ];
  }
