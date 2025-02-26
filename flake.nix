{
  description = "Description for the project";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";

    flake-parts = {
      url = "github:hercules-ci/flake-parts";
      inputs.nixpkgs-lib.follows = "nixpkgs";
    };

    flake-utils = {
      url = "github:numtide/flake-utils";
      inputs.systems.follows = "systems";
    };

    nix2container = {
      url = "github:nlewo/nix2container";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.flake-utils.follows = "flake-utils";
    };

    snowfall-lib = {
      url = "github:snowfallorg/lib";
      inputs.nixpkgs.follows = "nixpkgs";
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
    '';
    extra-trusted-public-keys = ''
      cache.nixos.org-1:6NCHdD59X431o0gWypbMrAURkbJ16ZPMQFGspcDShjY=
      nix-community.cachix.org-1:mB9FSh9qf2dCimDSUo8Zy7bkq5CX+/rkCWyvRCYg3Fs=
      devenv.cachix.org-1:w1cLUi8dv3hnoSPGAuibQv+f9TZLr6cv/Hm9XgU50cw=
    '';
    extra-experimental-features = "nix-command flakes";
  };

  outputs = inputs @ {flake-parts, ...}:
    flake-parts.lib.mkFlake {inherit inputs;} rec {
      systems = [
        "x86_64-linux"
        "i686-linux"
        "x86_64-darwin"
        "aarch64-linux"
        "aarch64-darwin"
      ];

      flake = let
        inherit (inputs) snowfall-lib;
        lib = snowfall-lib.mkLib {
          inherit inputs;
          src = ./.;

          snowfall = {
            namespace = "csnow";
            meta = {
              name = "csnow";
              title = "Conner Ohnesorge's Snowflake";
            };
          };
        };
      in
        lib.mkFlake {
          inherit inputs;

          outputs-builder = channels: {
            formatter = channels.nixpkgs.alejandra;
          };
        };
    };
}
