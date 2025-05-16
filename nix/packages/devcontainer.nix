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
  in {
    packages.x86_64-linux = rec {
      devcontainer = pkgs.dockerTools.streamNixShellImage {
        name = "conneroh/devcontainer";
        tag = "latest";
        drv = self.devShells.${system}.default;
      };

      deployDevcontainer = pkgs.writeShellApplication {
        name = "deploy-devcontainer";
        runtimeInputs = [
          pkgs.skopeo
        ];
        bashOptions = ["errexit" "pipefail"];
        text = ''
          set -e
          TOKEN=""

          [ -z "$GHCR_TOKEN" ] && GHCR_TOKEN="$(doppler secrets get --plain GHCR_TOKEN)"
          TOKEN="$GHCR_TOKEN"

          skopeo copy \
            --insecure-policy \
            docker-archive:"${devcontainer}" \
            "docker://$REGISTRY:latest" \
            --dest-creds x:"$TOKEN" \
            --format v2s2
        '';
      };
    };
  };
}
