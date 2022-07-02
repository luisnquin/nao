{
  buildGoPackage,
  fetchFromGitHub,
  lib,
}:
buildGoPackage rec {
  name = "nao";
  version = "example";
  rev = "";

  goPackagePath = "github.com/luisnquin/nao";

  src = fetchFromGitHub {
    inherit rev;
    owner = "luisnquin";
    repo = "nao";
    sha256 = null;
  };
  goDeps = ./deps.nix;

  ldflags = ["-s" "-w" "-X main.version=${version}"];

  meta = with lib; {
    # description = ""; # TODO
    license = licenses.mit;
    platforms = platforms.all;
    maintainers = with maintainers; [luisnquin];
  };
}
