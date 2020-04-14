GOARCH = amd64

UNAME = $(shell uname -s)

ifndef OS
	ifeq ($(UNAME), Linux)
		OS = linux
	else ifeq ($(UNAME), Darwin)
		OS = darwin
	endif
endif

#.DEFAULT_GOAL := all

# habit
.DEFAULT_GOAL := dev

all: fmt build start

dev:
	GOOS=$(OS) GOARCH="$(GOARCH)" go build -o vault-plugins-secrets-uuid cmd/uuid/main.go

build:
	GOOS=$(OS) GOARCH="$(GOARCH)" go build -o vault/plugins/uuid cmd/uuid/main.go

start:
	vault server -dev -dev-root-token-id=root -dev-plugin-dir=./vault/plugins

enable:
	vault secrets enable uuid

clean:
	rm -f ./vault/plugins/uuid

fmt:
	go fmt $$(go list ./...)

.PHONY: build clean fmt start enable
