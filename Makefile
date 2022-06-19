.PHONY: build
build:
	@go build -o ./build/nao ./src/main.go

run:
	@./build/nao