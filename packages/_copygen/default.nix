{
  lib,
  inputs,
  namespace,
  pkgs,
  stdenv,
  ...
}:
pkgs.buildGoModule rec {
  pname = "copygen";
  version = "0.4.1";

  src = pkgs.fetchFromGitHub {
    owner = "switchupcb";
    repo = "copygen";
    rev = "v${version}";
    sha256 = "sha256-gdoUvTla+fRoYayUeuRha8Dkix9ACxlt0tkac0CRqwA=";
  };

  vendorHash = "sha256-dOIGGZWtr8F82YJRXibdw3MvohLFBQxD+Y4OkZIJc2s=";
  subPackages = ["."];
  proxyVendor = true;

  ldflags = [
    "-s"
    "-w"
    "-X main.version=${version}"
  ];

  meta = with lib; {
    description = "HTTP linux cli proxy for monitoring and manipulating HTTP/HTTPS traffic";
    homepage = "https://github.com/monasticacademy/httptap";
    license = licenses.mit;
    maintainers = with maintainers; [
      connerohnesorge
      conneroisu
    ];
    mainProgram = "httptap";
  };
}
