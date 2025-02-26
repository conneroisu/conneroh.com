{
  # Snowfall Lib provides a customized `lib` instance with access to your flake's library
  # as well as the libraries available from your flake's inputs.
  lib,
  # You also have access to your flake's inputs.
  inputs,
  # The namespace used for your flake, defaulting to "internal" if not set.
  namespace,
  # All other arguments come from NixPkgs. You can use `pkgs` to pull packages or helpers
  # programmatically or you may add the named attributes as arguments here.
  pkgs,
  stdenv,
  ...
}:
pkgs.buildGo124Module {
  pname = "conneroh";
  version = "0.0.1";
  src = ./.;
  vendorSha256 = "sha256-0v/4+0/0/0";
  meta = {
    description = "A simple CLI tool to manage your connections";
    homepage = "https://github.com/conneroh/conneroh";
    license = lib.licenses.mit;
    maintainers = with lib.maintainers; [conneroh];
  };
}
