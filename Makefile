.PHONY: build
build:
	@go build -o build cmd/app/*.go

run:
	@./build

debug:
	@dlv debug cmd/app/main.go

all:
	@make build && make run
