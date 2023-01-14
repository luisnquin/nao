{}: let
  pkgs = import <nixpkgs> {};
in
  pkgs.buildGoModule rec {
    pname = "nao";
    version = "3.0.0";
    src = ./.;
    buildTarget = "./cmd/nao";

    vendorSha256 = "sha256-sDLbIHeY2Jdos8pgRaR6vbjwRPunAukasKz0m7sMIek=";

    doCheck = false;
    ldflags = ["-X main.version=${version}"];
  }
