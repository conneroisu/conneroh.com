{
  description = "Personal Website for Conner Ohnesorge";

  # TODO: Might be adventageous to introduce flake-utils to reduce copy-pasta of nix paths
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    systems.url = "github:nix-systems/default";
    bun2nix.url = "github:baileyluTCD/bun2nix";
    flake-parts.url = "github:hercules-ci/flake-parts";
    flake-parts.inputs.nixpkgs-lib.follows = "nixpkgs";
    rust-overlay.url = "github:oxalica/rust-overlay";
    rust-overlay.inputs.nixpkgs.follows = "nixpkgs";
    nix2container.url = "github:nlewo/nix2container";
    nix2container.inputs = {
      nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = inputs @ {flake-parts, ...}:
    flake-parts.lib.mkFlake {inherit inputs;} (
      {...}: {
        systems = import inputs.systems;
        imports = [
          ./nix/flake-module.nix
        ];

        perSystem = {system, ...}: {
          _module.args.pkgs = import inputs.nixpkgs {
            inherit system;
            config.allowUnfree = true;
            overlays = [
              (final: prev: {
                go = prev.go_1_24;
              })
              inputs.rust-overlay.overlays.default
            ];
          };
        };
      }
    );
}
