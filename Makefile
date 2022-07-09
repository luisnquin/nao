last_tag_released=$(shell git tag | tail -n 1)

.PHONY: build
build:
	@go build -ldflags "-s -w" -o ./build/nao .

clean:
	@rm ./build/*

run:
	@./build/nao

nix-build:
	@nix-build -E 'with import <nixpkgs> { };  callPackage ./default.nix {}' 

vue-dev: build
	@(cd client; npm run dev) & ./build/nao server -q -p=5000

sync:
	@bash ./tag-syncer.sh

install:
	@go install .

install-remote:
	@go install github.com/luisnquin/nao@$(last_tag_released)
