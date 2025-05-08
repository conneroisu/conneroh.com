{
  imports = [
    ./shells/flake-module.nix
  ];
  perSystem = {
    self',
    config,
    inputs',
    pkgs,
    system,
    ...
  }: {
  };
}
