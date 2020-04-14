GOARCH = amd64

UNAME = $(shell uname -s)
GO_CMD?=go
TEST?=$$($(GO_CMD) list ./... | grep -v /vendor/ | grep -v /integ)
TEST_TIMEOUT?=45m
EXTENDED_TEST_TIMEOUT=60m

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

test:
	@CGO_ENABLED=$(CGO_ENABLED) \
	VAULT_ADDR= \
	VAULT_TOKEN= \
	VAULT_DEV_ROOT_TOKEN_ID= \
	VAULT_ACC= \
	$(GO_CMD) test $(TEST) $(TESTARGS) -timeout=$(TEST_TIMEOUT) -parallel=20

# testacc runs acceptance tests
testacc: 
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package"; \
		exit 1; \
	fi
	VAULT_ACC=1 $(GO_CMD) test $(TEST) -v $(TESTARGS) -timeout=$(EXTENDED_TEST_TIMEOUT)

.PHONY: build clean fmt start enable
