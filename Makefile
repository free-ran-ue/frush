.PHONY: build

.DEFAULT_GOAL := build

build:
	go build -o frush main.go

clean:
	rm -rf frush