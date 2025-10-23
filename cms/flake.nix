{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    bun2nix.url = "github:baileyluTCD/bun2nix";
    bun2nix.inputs.nixpkgs.follows = "nixpkgs";
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
    bun2nix,
    ...
  } @ inputs:
    flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [inputs.hercules-ci-effects.overlay];
        };
      in {
        outputs = {
          default = let
            bunNix =
              pkgs.runCommand "bun.nix" {
                buildInputs = [bun2nix.packages.${system}.default];
              } ''
                bun2nix --lock-file ${./bun.lock} --output-file $out
              '';
          in
            bun2nix.lib.${system}.mkBunDerivation {
              src = self;
              bunLock = ./bun.lock;
              packageJson = ./package.json;
              inherit bunNix;

              preBuild = ''
                bun run build
                mkdir -p $out/bin/dist
                cp -r dist/* $out/bin/dist
                mkdir -p $out/bin/node_modules
                cp -r node_modules/* $out/bin/node_modules
              '';

              buildFlags = [
                "--compile"
                "--minify"
                "--sourcemap"
              ];

              postInstall = ''
                # Create wrapper script to set Prisma engine paths for NixOS
                mv $out/bin/cms $out/bin/.connix-unwrapped

                cat > $out/bin/cms <<EOF
                  #!${pkgs.bash}/bin/bash
                  # Execute the wrapped binary using self-referential path
                  SCRIPT_DIR="\$(cd "\$(dirname "\$0")" && pwd)"
                  exec "\$SCRIPT_DIR/.cms-unwrapped" "\$@"
                EOF

                chmod +x $out/bin/cms
              '';

              nativeBuildInputs = with pkgs; [
                openssl.dev
                openssl
              ];

              meta = with pkgs.lib; {
                description = "TanStack Start fullstack application";
                license = licenses.mit;
              };
            };
        };
      }
    );
}
