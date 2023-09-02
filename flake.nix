{
  description = "A CLI tool to take notes without worrying about the path where the file is";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = inputs:
    with inputs;
      flake-utils.lib.eachDefaultSystem (
        system: let
          pkgs = import nixpkgs {
            inherit system;
          };

          defaultPackage = pkgs.callPackage ./default.nix {};
        in {
          inherit defaultPackage;

          defaultApp = flake-utils.lib.mkApp {
            drv = defaultPackage;
          };

          devShell = pkgs.mkShell {
            buildInputs = [defaultPackage];
          };
        }
      );
}
