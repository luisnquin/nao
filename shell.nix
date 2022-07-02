with import <nixpkgs> {};
  pkgs.mkShell {
    buildInput = [pkgs.go pkgs.gcc];
  }
