.PHONY: build
build:
	@go build  -o ./build/nao ./cmd/nao/

clean:
	@rm ./build/*

run:
	@./build/nao 

sync:
	@bash ./tools/tag-syncer.sh

install:
	@go install ./cmd/nao/

uninstall:
	@rm -f ~/go/bin/nao

nix-install:
	@nix-build default.nix

nix-uninstall:
	@nix-env -e nao

nix-clean:
	@ls /nix/store | grep nao | xargs -I {} bash -c 'nix-store --delete /nix/store/{}'

style: 
	@gofumpt -w ./internal