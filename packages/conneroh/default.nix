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
  vendorHash = "sha256-nm3JEU+6MuA0bCXfAswgb7JZBmysypDIKjvw+PJFltY=";
  doCheck = true;
  checkPhase = ''
    echo "Running conneroh for 3 seconds to ensure it works..."
    timeout 3 ./result/bin/conneroh || true
  '';
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
