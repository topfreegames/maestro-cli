# maestro-cli api
# https://github.com/topfreegames/mystack-controller
#
# Licensed under the MIT license:
# http://www.opensource.org/licenses/mit-license
# Copyright © 2017 Top Free Games <backend@tfgco.com>

TEST_PACKAGES=`find . -type f -name "*.go" ! \( -path "*vendor*" \) | sed -En "s/([^\.])\/.*/\1/p" | uniq`
BIN_PATH = "./bin"
BIN_NAME = "maestro"

build:
	@mkdir -p bin
	@go build -o ./bin/maestro main.go

test: unit test-coverage-func
	
unit: clear-coverage-profiles unit-run gather-unit-profiles

merge-profiles:
	@mkdir -p _build
	@gocovmerge _build/*.out > _build/coverage-all.out

test-coverage-func coverage-func: merge-profiles
	@echo
	@echo "=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-"
	@echo "Functions NOT COVERED by Tests"
	@echo "=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-"
	@go tool cover -func=_build/coverage-all.out | egrep -v "100.0[%]"

clear-coverage-profiles:
	@find . -name '*.coverprofile' -delete

unit-run:
	@ginkgo -tags unit -cover -r -randomizeAllSpecs -randomizeSuites -skipMeasurements ${TEST_PACKAGES}

gather-unit-profiles:
	@mkdir -p _build
	@echo "mode: count" > _build/coverage-unit.out
	@bash -c 'for f in $$(find . -name "*.coverprofile"); do tail -n +2 $$f >> _build/coverage-unit.out; done'

build-all-platforms:
	@mkdir -p ${BIN_PATH}
	@echo "Building for linux-i386..."
	@env GOOS=linux GOARCH=386 go build -o ${BIN_PATH}/${BIN_NAME}-linux-i386
	@echo "Building for linux-x86_64..."
	@env GOOS=linux GOARCH=amd64 go build -o ${BIN_PATH}/${BIN_NAME}-linux-x86_64
	@echo "Building for darwin-i386..."
	@env GOOS=darwin GOARCH=386 go build -o ${BIN_PATH}/${BIN_NAME}-darwin-i386
	@echo "Building for darwin-x86_64..."
	@env GOOS=darwin GOARCH=amd64 go build -o ${BIN_PATH}/${BIN_NAME}-darwin-x86_64
	@echo "Building for win-x86_64..."
	@env GOOS=windows GOARCH=amd64 go build -o ${BIN_PATH}/${BIN_NAME}-win-x86_64
