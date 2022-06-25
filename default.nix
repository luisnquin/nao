{ buildGoPackage, fetchgit, lib }:

buildGoPackage rec {
  name = "vgo2nix-${version}";
  version = "example";
  # rev = "";

  goPackagePath = "github.com/luisnquin/nao";

  src = fetchgit {
    # inherit rev;
    url = "git@github.com:luisnquin/nao.git";
    sha256 = null;
  };
  goDeps = ./deps.nix;

  meta = with lib; {
    # description = ""; # TODO
    license = licenses.mit;
    platforms = platforms.all;
    maintainers = with maintainers; [ luisnquin ];
  };
}
