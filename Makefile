# maestro-cli api
# https://github.com/topfreegames/maestro-cli
#
# Licensed under the MIT license:
# http://www.opensource.org/licenses/mit-license
# Copyright © 2017 Top Free Games <backend@tfgco.com>

.PHONY: mocks

TEST_PACKAGES=`find . -type f -name "*.go" ! \( -path "*vendor*" \) | sed -En "s/([^\.])\/.*/\1/p" | uniq`
BIN_PATH = "./bin"
BIN_NAME = "maestro"

SOURCES := $(shell \
	find . -not \( \( -name .git -o -name .go -o -name vendor \) -prune \) \
	-name '*.go')

.PHONY: help
help: Makefile ## Show list of commands.
	@echo "Choose a command to run in "$(APP_NAME)":"
	@echo ""
	@awk 'BEGIN {FS = ":.*?## "} /[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-40s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort


build: ## Generate cli binary from code.
	@mkdir -p bin
	@go build -o ./bin/maestro-cli main.go

test: unit test-coverage-func ## Run coverage and unit tests.
	
unit: clear-coverage-profiles unit-run gather-unit-profiles ## Run unit tests.

merge-profiles: ## Merge coverage profiles
	@mkdir -p _build
	@go run github.com/wadey/gocovmerge _build/*.out > _build/coverage-all.out

test-coverage-func coverage-func: merge-profiles ## Execute tests coverage.
	@echo
	@echo "=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-"
	@echo "Functions NOT COVERED by Tests"
	@echo "=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-"
	@go tool cover -func=_build/coverage-all.out | egrep -v "100.0[%]"

.PHONY: clear-coverage-profiles
clear-coverage-profiles: ## Delete coverage profiles.
	@find . -name '*.coverprofile' -delete

.PHONY: unit-run
unit-run: ## Execute unit tests.
	@go run github.com/onsi/ginkgo/ginkgo -tags unit -cover -r -randomizeAllSpecs -randomizeSuites -skipMeasurements ${TEST_PACKAGES}

.PHONY: gather-unit-profiles
gather-unit-profiles: ## Gather unit profiles.
	@mkdir -p _build
	@echo "mode: count" > _build/coverage-unit.out
	@bash -c 'for f in $$(find . -name "*.coverprofile"); do tail -n +2 $$f >> _build/coverage-unit.out; done'

.PHONY: build-all-platforms
build-all-platforms: ## Build binaries for all supported platforms.
	@mkdir -p ${BIN_PATH}
	@echo "Building for linux-i386..."
	@env GOOS=linux GOARCH=386 go build -o ${BIN_PATH}/${BIN_NAME}-linux-i386
	@echo "Building for linux-x86_64..."
	@env GOOS=linux GOARCH=amd64 go build -o ${BIN_PATH}/${BIN_NAME}-linux-x86_64
	@echo "Building for darwin-x86_64..."
	@env GOOS=darwin GOARCH=amd64 go build -o ${BIN_PATH}/${BIN_NAME}-darwin-x86_64
	@echo "Building for win-x86_64..."
	@env GOOS=windows GOARCH=amd64 go build -o ${BIN_PATH}/${BIN_NAME}-win-x86_64

.PHONY: mocks
mocks: ## Generate mocks.
	@echo 'making mocks from ./interfaces'
	mockgen -source=interfaces/client.go -destination=mocks/client.go -package=mocks
	mockgen -source=interfaces/filesystem.go -destination=mocks/filesystem.go -package=mocks
	@echo 'done, mocks on ./mocks'

.PHONY: goimports
goimports: ## Run goimports to format files.
	@go run golang.org/x/tools/cmd/goimports -w $(SOURCES)

.PHONY: lint
lint: ## Run linter.
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint run -E goimports ./...
