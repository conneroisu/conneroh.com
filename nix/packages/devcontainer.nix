{
  inputs,
  self,
  ...
}: let
  system = "x86_64-linux";
in {
  flake = let
    pkgs = import inputs.nixpkgs {
      inherit system;
    };
    tag = "v5";
  in {
    packages.x86_64-linux = rec {
      devcontainer = pkgs.dockerTools.buildNixShellImage {
        inherit tag;
        name = "conneroh/devcontainer";
        drv = self.devShells.${system}.devcontainer;
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
          REGISTRY="ghcr.io/conneroisu/conneroh.com"

          skopeo copy \
            --insecure-policy \
            docker-archive:"${devcontainer}" \
            "docker://$REGISTRY:${tag}" \
            --dest-creds x:"$TOKEN" \
            --format v2s2
        '';
      };
    };
  };
}
