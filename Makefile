# (C) Copyright 2022 Hewlett Packard Enterprise Development LP

# ===================================================
# Makefile for building Panorama Fleet Hauler service
# ===================================================
#
# Configuration
# =============
# These can be overwritten on invocation (e.g `make PREFIX=/usr install`)
#

SHELL=/bin/bash -o pipefail

DOCKER   := docker
GO       := go
KUBECTL  := kubectl
export GO111MODULE := auto
export GOFLAGS     := -mod=vendor
export GOPRIVATE   := github.hpe.com

GO_BUILD := $(GO) build
GOTEST  := $(GO) test

SERVER_NAME := panorama-common-hauler
DIST_DIR := ./dist
CMD_DIR  := ./cmd
INTERNAL_DIR := ./internal
FUNCTIONAL_TEST_DIR :=./tests/component
CONTRACT_CONSUMER_TEST_DIR :=./tests/contract/consumer
CONTRACT_PROVIDER_TEST_DIR :=./tests/contract/provider
ifndef BRANCH_NAME
export BRANCH_NAME := $(shell git rev-parse --abbrev-ref HEAD)
endif
AUTHZ_MOCK_SERVER :=./tests/mock-servers/authz-mock
FLEET_MOCK_SERVER :=./tests/mock-servers/fleet-mock

VERSION = $(shell (git describe --long --tags --match 'v[0-9]*' || echo v0.0.0) | cut -c2-)
LDFLAGS = -X main.Version=$(VERSION)
TEST_MODULES = $(shell $(GO) list $(INTERNAL_DIR)/...)
FUNCTIONAL_TEST_MODULES = $(shell $(GO) list $(FUNCTIONAL_TEST_DIR)/...)
CONTRACT_CONSUMER_TEST_MODULES_MAIN = $(shell $(GO) list $(CONTRACT_CONSUMER_TEST_DIR)/...)
CONTRACT_PROVIDER_TEST_MODULES_MAIN = $(shell $(GO) list $(CONTRACT_PROVIDER_TEST_DIR)/...)

COVER_PROFILE := $(SERVER_NAME).cov

# The CI_MINIMUM_TEST_COVERAGE environment variable is set automatically in
# Jenkins runs.  A default value is defined here for developer builds.
# Will continue working on Unit tests to improve the code coverage.
CI_MINIMUM_TEST_COVERAGE = 85

#
# Targets
# =======

## help: Output this message and exit.
help:
	@fgrep -h '##' $(MAKEFILE_LIST) | fgrep -v fgrep | column -t -s ':' | sed -e 's/## //'
.PHONY: help

## all: lint and build
.PHONY: all
all: lint build

## build: build the project
.PHONY: build
build: vendor #generate_mocks
	# Since Mac uses nfs to mount the repositories, the user gets squashed down
	# to the underlying Mac user causing permission issues. The workaround is to
	# create a directory here which will be done within the docker container's overlay
	# filesystem as the Dockerfile calls the Makefile
	@mkdir -p dist
	$(GO_BUILD) -ldflags "$(LDFLAGS)" -o dist/ ./cmd/...

## vendor: download vendored dependencies
vendor: go.mod go.sum
	$(GO) mod vendor

## clean: clean up built code and vendor directory
.PHONY: clean
clean:
	rm -rf c.out coverage.html json-report.txt component.test $(DIST_DIR) $(COVER_PROFILE) consumer-contract-json-report.txt provider-contract-json-report.txt allure-reports vendor

.PHONY: golanglint
golanglint: vendor
	golangci-lint cache clean
	golangci-lint -j 2 run --timeout 10m

## lint: run the project linters. Project linters configured in .golangci.yml in project root.
.PHONY: lint
lint: golanglint

## component-test: Build a component test executable, to be executed against a running panorama-fleet-hauler service and related infrastructure.
component-test: vendor
	$(GOTEST) -c -v $(FUNCTIONAL_TEST_MODULES)
.PHONY: component-test

## unit-test: Run all the unit tests
.PHONY: unit-test
unit-test: vendor
	GOPROXY='' GOFLAGS='' $(GO) install github.com/ilyubin/gotest2allure/cmd/gotest2allure@latest
	$(GO) mod vendor
	$(GOTEST) $(TEST_MODULES) \
		-v \
		-cover \
		-coverprofile=c.out \
		-count=1 \
		-json \
		> json-report.txt || (jq . < json-report.txt; false)
	cat c.out | grep -v "/external/" | grep -v "/mocks/" > $(COVER_PROFILE)
	$(HOME)/go/bin/gotest2allure -f json-report.txt -o allure-reports/unit-test-reports

## code-coverage: Produce a per-function code coverage report from a unit test run.
.PHONY: code-coverage
code-coverage: $(COVER_PROFILE)
	$(GO) tool cover -func $(COVER_PROFILE)
	$(GO) tool cover -html=$(COVER_PROFILE) -o unit_test_coverage.html
	$(eval TEST_COVERAGE = $(shell $(GO) tool cover -func $(COVER_PROFILE) | grep 'total:' | awk '{print substr($$3, 1, length($$3)-1)}'))
	@echo "Unit-tests passed with $(TEST_COVERAGE)  coverage"
	@if [ `expr $(TEST_COVERAGE) \< $(CI_MINIMUM_TEST_COVERAGE)` -eq 1 ]; then echo "Require atleast ${CI_MINIMUM_TEST_COVERAGE}% test coverage"; exit 1; fi

## consumer-contract-test: Run consumer contract tests
.PHONY: consumer-contract-test
consumer-contract-test: vendor
	GOPROXY='' GOFLAGS='' $(GO) install github.com/ilyubin/gotest2allure/cmd/gotest2allure@latest
	$(GO) mod vendor
	$(GOTEST) $(CONTRACT_CONSUMER_TEST_MODULES_MAIN) \
		-v \
		-cover \
		-coverprofile=c.out \
		-count=1 -json > consumer-contract-json-report.txt
		$(HOME)/go/bin/gotest2allure -f consumer-contract-json-report.txt -o allure-reports/consumer-contract-test-reports

## publish-pact: Run publish pact
.PHONY: publish-pact
publish-pact:
	@echo "Publishing Pact file..."
	pact-broker publish ./tests/contract/pacts --consumer-app-version="${PACT_VERSION}" --tag ${BRANCH_NAME} --branch ${BRANCH_NAME} --broker-username ${PACT_USER} --broker-password ${PACT_PASS} --broker-base-url ${PACT_BROKER_URL}

## provider-contract-test: Run provider contract tests
.PHONY: provider-contract-test
provider-contract-test: vendor
	GOPROXY='' GOFLAGS='' $(GO) install github.com/ilyubin/gotest2allure/cmd/gotest2allure@latest
	$(GO) mod vendor
	$(GOTEST) $(CONTRACT_PROVIDER_TEST_MODULES_MAIN) \
		-v \
		-cover \
		-coverprofile=c.out \
		-count=1 -json > provider-contract-json-report.txt
		$(HOME)/go/bin/gotest2allure -f provider-contract-json-report.txt -o allure-reports/provider-contract-test-reports

## can-i-deploy: verifies pacts from pact-broker
.PHONY: can-i-deploy
can-i-deploy:
	@echo "Running can-i-deploy to check if panorama-fleet-hauler can be deployed."
	pact-broker can-i-deploy --pacticipant panorama-fleet-hauler \
	--retry-while-unknown=2 --retry-interval=30 \
	--version $(PACT_VERSION) --broker-base-url="$(PACT_BROKER_URL)" \
	--broker-username="$(PACT_USER)" --broker-password="$(PACT_PASS)"

.PHONY: setup_kafka
setup_kafka:
	echo NotImplemented yet

## setup_env: Setup the environment - assumes the ccs-* services are running
.PHONY: setup_env
setup_env: setup_kafka

## mock-servers-image: Build all mock server images - Authz, Fleet, HCM, VirtAggr.
mock-servers-image: vendor
	$(DOCKER) build -t mock-authz:1.0 -f $(AUTHZ_MOCK_SERVER)/Dockerfile .
	$(DOCKER) build -t mock-fleet:1.0 -f $(FLEET_MOCK_SERVER)/Dockerfile .
.PHONY: mock-servers-image

## mock-servers-deploy: Deploy all mock servers.
mock-servers-deploy:
	$(KUBECTL) apply -f $(AUTHZ_MOCK_SERVER)
	$(KUBECTL) apply -f $(FLEET_MOCK_SERVER)
.PHONY: mock-servers-deploy

## mock-servers-undeploy: Undeploy all mock servers.
mock-servers-undeploy:
	$(KUBECTL) delete -f $(AUTHZ_MOCK_SERVER)
	$(KUBECTL) delete -f $(FLEET_MOCK_SERVER)
.PHONY: mock-servers-undeploy

## mock: Create mockery mock for specified interface in ./internal/mocks directory following mock file/struct name conventions, e.g. make mock interface=DatastoreServiceInterface packagedir=./internal/services
.PHONY: mock
mock:
	mockery --name=$(interface) \
                --dir=$(packagedir) \
                --output ./internal/mocks \
                --filename mock_$(interface).go \
                --structname Mock$(interface) \
                --keeptree