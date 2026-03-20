{
  description = "A CLI tool to take notes without worrying about the path where the file is";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    systems.url = "github:nix-systems/default-linux";
  };

  outputs = {
    self,
    nixpkgs,
    systems,
    ...
  }: let
    inherit (nixpkgs) lib;
    eachSystem = lib.genAttrs (import systems);
    pkgsFor = eachSystem (system:
      import nixpkgs {
        localSystem = system;
      });
  in {
    packages = eachSystem (system: rec {
      nao = pkgsFor.${system}.callPackage ./default.nix {};
      default = nao;
    });

    overlays.default = final: prev: {
      nao = self.packages.${final.system}.default;
    };

    homeManagerModules.default = import ./nix/hm-module.nix self;
  };
}
