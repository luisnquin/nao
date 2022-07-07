last_tag_released=$(shell git tag | tail -n 1)

.PHONY: build
build:
	@go build -ldflags "-s -w" -o ./build/nao ./src/cmd/nao/main.go

run:
	@./build/nao

nix-build:
	@nix-build -E 'with import <nixpkgs> { };  callPackage ./default.nix {}' 

vue-dev: build
	@(cd client; npm run dev) & ./build/nao server -q -p=5000

sync:
	@bash ./tag-syncer.sh

install:
	@go install ./src/cmd/nao

install-remote:
	@go install github.com/luisnquin/nao/cmd/nao@$(last_tag_released)
