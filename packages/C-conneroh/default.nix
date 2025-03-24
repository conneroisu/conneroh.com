{
  lib,
  inputs,
  namespace,
  pkgs,
  stdenv,
  ...
}:
pkgs.dockerTools.buildLayeredImage {
  name = "conneroh.com";
  tag = "latest";
  contents = [
    pkgs."${namespace}".conneroh
    pkgs.cacert
  ];
  config = {
    WorkingDir = "/root";
    Cmd = ["/bin/conneroh.com"];
    ExposedPorts = {
      "8080/tcp" = {};
    };
    Env = [
      "SSL_CERT_FILE=${pkgs.cacert}/etc/ssl/certs/ca-bundle.crt"
      "NIX_SSL_CERT_FILE=${pkgs.cacert}/etc/ssl/certs/ca-bundle.crt"
    ];
  };
  extraCommands = ''
    echo "$(git rev-parse HEAD)" > REVISION
  '';
}
