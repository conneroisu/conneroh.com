{
  lib,
  inputs,
  namespace,
  pkgs,
  stdenv,
  ...
}:
pkgs.buildGo124Module {
  pname = "conneroh";
  version = "0.0.1";
  src = ./../../.;
  subPackages = ["cmd/conneroh"];
  vendorSha256 = "";
  meta = {
    description = "Personal Website";
    homepage = "https://github.com/conneroisu/conneroh.com";
    license = lib.licenses.mit;
    maintainers = with lib.maintainers; [conneroh];
  };
}
