{
  inputs.conneroh.url = "github:conneroisu/conneroh";
  inputs.nix2container.url = "github:nlewo/nix2container";
  outputs = { conneroh, nix2container, ... }:
  {
    nix2container.buildImage {
      name = "layered";
      config = {
        entrypoint = ["${pkgs.hello}/bin/hello"];
      };
      maxLayers = 3;
      fromImage = nix2container.pullImage {
        imageName = "alpine";
        imageDigest = "sha256:115731bab0862031b44766733890091c17924f9b7781b79997f5f163be262178";
        arch = "amd64";
        sha256 = "sha256-o4GvFCq6pvzASvlI5BLnk+Y4UN6qKL2dowuT0cp8q7Q=";
      };
      metadata = {
        created_by = "test created_by";
        author = "test author";
        comment = "test comment";
      };
    }
  }
}