{
  perSystem = {
    self',
    config,
    inputs',
    pkgs,
    ...
  }: let
    buildWithSpecificGo = pkg: pkg.override {buildGoModule = pkgs.buildGo124Module;};
    shell-scripts = import ./shell-scripts.nix {
      inherit pkgs self';
    };
    inherit (shell-scripts) scripts scriptPackages;
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

          litellm

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
          sqlite

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

    packages = import ./shell-packages.nix {
      inherit self' pkgs scripts scriptPackages;
    };
  };
}
