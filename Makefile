.PHONY: build
build:
	@go build -ldflags "-s -w" -o ./build/nao ./main.go

run:
	@./build/nao

nix-build:
	@nix-build -E 'with import <nixpkgs> { };  callPackage ./default.nix {}' 