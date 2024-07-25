.PHONY: default
default: build

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

fname := app
fext :=
# ifeq ($(GOOS),Windows_NT)
ifeq ($(GOOS),windows)
	fext := .exe
endif


build:
#	go build -trimpath -ldflags="-w -s" -o bin/$(fname)$(fext) ./cmd
	go build -trimpath -ldflags="-w -s" -o bin/$(fname)$(fext) .

clean:
	rm -rf ./bin/*
