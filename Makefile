.PHONY: build

.DEFAULT_GOAL := build

build:
	go build -o bin/frush main.go

clean:
	rm -rf bin