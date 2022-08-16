GOCMD=go
BINARY_NAME=tokencli
CONTRACT_NAME=ERC20Token
EXPORT_RESULT?=false

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

.PHONY: build

## Build:
build: ## Build your project and put the output binary in out/bin/
	mkdir -p out/bin
	GO111MODULE=on $(GOCMD) build -mod vendor -o out/bin/$(BINARY_NAME) ./cmd/tokencli

install:
	go install ./cmd/tokencli

vendor: ## Copy of all packages needed to support builds and tests in the vendor directory
	$(GOCMD) mod vendor

clean: ## Remove all build related files
	rm -fr out
	rm -f checkstyle-report.xml

## Compile:
compile: ## Compile Solidity files and generate Go bindings
	solc --optimize --abi ./contracts/$(CONTRACT_NAME).sol -o build --include-path node_modules/ --base-path .
	solc --optimize --bin ./contracts/$(CONTRACT_NAME).sol -o build --include-path node_modules/ --base-path .
	abigen --abi=./build/$(CONTRACT_NAME).abi --bin=./build/$(CONTRACT_NAME).bin --pkg=token --out=./contracts/ERC20Token.go

clean-sol: ## Clean all Solidity compiled files and Go bindings
	rm -fr build/*
	rm contracts/*.go

## Lint:
lint: ## Run golintci-lint
	$(eval OUTPUT_OPTIONS = $(shell [ "${EXPORT_RESULT}" == "true" ] && echo "--out-format checkstyle ./... | tee /dev/tty > checkstyle-report.xml" || echo "" ))
	golangci-lint run --deadline=65s $(OUTPUT_OPTIONS)

## Help:
help: ## Show this help.
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)