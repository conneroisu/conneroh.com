{
  inputs,
  self,
  ...
}: let
  system = "x86_64-linux";
in {
  flake = let
    inherit (inputs.nix2container.packages.x86_64-linux) nix2container;
    pkgs = import inputs.nixpkgs {
      inherit system;
    };
    alpine = nix2container.pullImage {
      imageName = "alpine";
      imageDigest = "sha256:115731bab0862031b44766733890091c17924f9b7781b79997f5f163be262178";
      arch = "amd64";
      sha256 = "sha256-o4GvFCq6pvzASvlI5BLnk+Y4UN6qKL2dowuT0cp8q7Q=";
    };
  in {
    packages.x86_64-linux.devcontainer = nix2container.buildImage {
      name = "devcontainer";
      initializeNixDatabase = true;
      fromImage = alpine;
      layers = [
        (nix2container.buildLayer {
          deps = [
            pkgs.direnv
          ];
        })
        (nix2container.buildLayer {
          deps = [
            self.devShells.${system}.default
          ];
        })
      ];
      config = {
        Env = [
          "NIX_PAGER=cat"
          # A user is required by nix
          # https://github.com/NixOS/nix/blob/9348f9291e5d9e4ba3c4347ea1b235640f54fd79/src/libutil/util.cc#L478
          "USER=nobody"
          # When running in podman on the GitHub CI, Nix fails to find the
          # user home dir for an unkonwn reason...
          "HOME=/"
        ];
      };
    };
  };
}
