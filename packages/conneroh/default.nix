{
  lib,
  inputs,
  namespace,
  pkgs,
  stdenv,
  ...
}:
pkgs.buildGoModule {
  pname = "conneroh";
  version = "0.0.1";
  src = ../../.;
  subPackages = ["."];
  nativeBuildInputs = [];
  vendorHash = "sha256-nm3JEU+6MuA0bCXfAswgb7JZBmysypDIKjvw+PJFltY=";
  preBuild = ''
    ${pkgs.templ}/bin/templ generate
  '';
  meta = {
    description = "Personal Website";
    homepage = "https://github.com/conneroisu/conneroh.com";
    license = lib.licenses.mit;
    maintainers = with lib.maintainers; [conneroh];
  };
}
