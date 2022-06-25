.PHONY: build
build:
	@go build -ldflags "-s -w" -o ./build/nao ./src/main.go

run:
	@./build/nao