{
  # Snowfall Lib provides a customized `lib` instance with access to your flake's library
  # as well as the libraries available from your flake's inputs.
  lib,
  # You also have access to your flake's inputs.
  inputs,
  # The namespace used for your flake, defaulting to "internal" if not set.
  namespace,
  # All other arguments come from NixPkgs. You can use `pkgs` to pull checks or helpers
  # programmatically or you may add the named attributes as arguments here.
  pkgs,
  ...
}:
inputs.pre-commit-hooks.lib.${pkgs.system}.run {
  src = ./../../.;
  hooks = {
    alejandra.enable = true;
    gofmt.enable = true;
    prettier = {
      enable = true;
      excludes = [
        "flake.lock"
        ".*_dist/"
        ".*node_modules/"
      ];
      args = [
        "--ignore-path=.gitignore"
        "--log-level=debug"
      ];
    };
  };
}
