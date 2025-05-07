{
  description = "Personal Website Production Testing";

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
    in rec {
      packages = {
        live-ci = pkgs.buildGoModule {
          src = ./../../.;
          vendorHash = null;
          version = self.shortRev or "dirty";
          name = "live-ci";
          subPackages = ["./cmd/live-ci"];
        };
        run = pkgs.writeShellScriptBin "run" {
          text = ''
            export REPO_ROOT=$(git rev-parse --show-toplevel)
            export PLAYWRIGHT_BROWSERS_PATH=${pkgs.playwright-driver.browsers}
            export PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD=1
            export PLAYWRIGHT_NODEJS_PATH=${pkgs.nodejs_20}/bin/node

            # Browser executable paths
            export PLAYWRIGHT_CHROMIUM_EXECUTABLE_PATH=${ #
              "${pkgs.playwright-driver.browsers}/chromium-1155"
            }
            ${packages.live-ci}/bin/live-ci
          '';
        };
      };
    });
}
