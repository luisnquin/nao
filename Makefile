last_tag_released=$(shell git tag | tail -n 1)

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

install-remote:
	@go install github.com/luisnquin/nao/cmd/nao/@$(last_tag_released)

style: 
	@gofumpt -w ./