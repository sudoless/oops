#      @ SUDOLESS SRL <contact@sudoless.org>
#      This Source Code Form is subject to the
#      terms of the Mozilla Public License, v.
#      2.0. If a copy of the MPL was not
#      distributed with this file, You can
#      obtain one at
#      http://mozilla.org/MPL/2.0/.


THIS_MAKEFILE_VERSION = v0.2.1
THIS_MAKEFILE_UPDATE = master
THIS_MAKEFILE := $(lastword $(MAKEFILE_LIST))
THIS_MAKEFILE_URL := https://raw.githubusercontent.com/sudoless/make/$(THIS_MAKEFILE_UPDATE)/go.mk


# GO
export CGO_ENABLED ?= 0
GO ?= GO111MODULE=on go
GO_TAGS ?= timetzdata

# SOURCE
SOURCE_FILES?=$$(find . -name '*.go' | grep -v pb.go | grep -v vendor)

# OUTPUT
DIR_OUT   := out
FILE_COV  := $(DIR_OUT)/cover.out

# MOD
ifneq ("$(wildcard go.mod/)","") # check go.mod exists
PROJECT_MOD_NAME := $(shell go list -m -mod=readonly)
endif


# DEV DEPS
DEV_EXTERNAL_TOOLS=\
	github.com/golangci/golangci-lint/cmd/golangci-lint@v1.43.0 \
	github.com/securego/gosec/v2/cmd/gosec@v2.9.5 \
	github.com/fzipp/gocyclo/cmd/gocyclo@v0.4.0 \
	github.com/jstemmer/go-junit-report@v0.9.1 \
	golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@v0.1.8 \
	mvdan.cc/gofumpt@v0.2.1 \
	gotest.tools/gotestsum@v1.7.0

# OTHER
SHA256 ?= shasum -a 256


$(DIR_OUT):
	@mkdir -p $(DIR_OUT)

.PHONY: init
init: ## setup a barebones Go project
	@mkdir -p cmd/ pkg/ docs/ scripts/ data/ deployment/ internal/
	@echo "out/" >> .gitignore
	@echo "vendor/" >> .gitignore

.PHONY: run/%
run/%: build/% ## run the specified target
	@printf "$(FMT_PRFX) running $(FMT_INFO)$*$(FMT_END) from $(FMT_INFO)$(DIR_OUT)/dist/$*_$$(go env GOOS)_$$(go env GOARCH)$(FMT_END)\n"
	@$(DIR_OUT)/dist/$*_$$(go env GOOS)_$$(go env GOARCH)

.PHONY: build/%
build/%: APP_OUT ?= $*_$$(go env GOOS)_$$(go env GOARCH)
build/%: ## build a specific cmd/$(TARGET)/... into $(DIR_OUT)/dist/$(TARGET)...
	@printf "$(FMT_PRFX) building $(FMT_INFO)$*$(FMT_END) version=$(FMT_INFO)$(BUILD_VERSION)$(FMT_END)\
 buildhash=$(FMT_INFO)$(BUILD_HASH)$(FMT_END)\n"
	@printf "$(FMT_PRFX) using $(FMT_INFO)$$(go version)$(FMT_END)\n"
	@$(GO) build -trimpath -tags "$(GO_TAGS)" \
		-ldflags="-w -s \
			-X main._serviceName=$*           \
			-X main._version=$(BUILD_VERSION) \
			-X main._buildTime=$(BUILD_TIME)  \
			-X main._buildHash=$(BUILD_HASH)" \
		-o $(DIR_OUT)/dist/$(APP_OUT) \
		./cmd/$*/...
	@printf "$(FMT_PRFX) built binary $(FMT_INFO)$(DIR_OUT)/dist/$(APP_OUT)$(FMT_END)\n"

.PHONY: install/%
install/%: APP_OUT ?= $*_$$(go env GOOS)_$$(go env GOARCH)
install/%: GOBIN ?= $(GOPATH)/bin
install/%: build/% ## install the built binary to $GOBIN
	@printf "$(FMT_PRFX) installing $(FMT_INFO)$*$(FMT_END) ($(FMT_WARN)$(APP_OUT)$(FMT_END)) at $(FMT_INFO)$(GOBIN)/$*$(FMT_END)\n"
	@cp $(DIR_OUT)/dist/$(APP_OUT) $(GOBIN)/$*
	@printf "$(FMT_PRFX) $(FMT_OK)ok$(FMT_END) (which=$(FMT_INFO)$(shell which $*)$(FMT_END))\n"

.PHONY: clean
clean: ## remove build time generated files
	@printf "$(FMT_PRFX) removing output directory\n"
	@rm -rf $(DIR_OUT)/

.PHONY: purge
purge: clean ## remove everything that could cause environment issues
	@printf "$(FMT_PRFX) deleting system32\n"
	$(GO) clean -cache
	$(GO) clean -testcache
	$(GO) clean -modcache


.PHONY: test
test: export CGO_ENABLED=1
test: $(DIR_OUT) ## run unit tests
	@printf "$(FMT_PRFX) running tests\n"
	@gotestsum \
		--junitfile $(FILE_COV).xml \
		--format short -- \
		-race \
		-timeout=30s -parallel=20 -failfast \
		-covermode=atomic -coverpkg=./... -coverprofile=$(FILE_COV).txt \
		./...

.PHONY: test-deps
test-deps: ## run tests with dependencies
	@printf "$(FMT_PRFX) running all tests\n"
	$(GO) test all

.PHONY: bench
bench: ## run benchmarks
	@printf "$(FMT_PRFX) running benchmarks\n"
	$(GO) test -exclude-dir=vendor -exclude-dir=.cache -bench=. -benchmem -benchtime=10s ./...

.PHONY: cover
cover: ## open coverage file in browser
	@printf "$(FMT_PRFX) opening coverage file in browser\n"
	$(GO) tool cover -html=$(FILE_COV).txt

.PHONY: tidy
tidy: ## tidy and verify go modules
	@printf "$(FMT_PRFX) tidying go modules\n"
	$(GO) mod tidy
	$(GO) mod verify

.PHONY: download
download: ## download go modules
	@printf "$(FMT_PRFX) downloading dependencies as modules\n"
	@$(GO) mod $(GO_MOD) download -x

.PHONY: vendor
vendor: ## tidy, vendor and verify dependencies
	@printf "$(FMT_PRFX) downloading and creating vendor dependencies\n"
	$(GO) mod tidy -v
	$(GO) mod vendor -v
	$(GO) mod verify

.PHONY: updates
updates: ## display outdated direct dependencies
	@printf "$(FMT_PRFX) checking for direct dependencies updates\n"
	@$(GO) list -u -m -mod=readonly -json all | go-mod-outdated -direct

.PHONY: lint
lint: ## run golangci linter
	@printf "$(FMT_PRFX) running golangci-lint\n"
	@golangci-lint run -v --timeout 10m --skip-dirs=".cache/|vendor/|scripts/|docs/|deployment/|data/"  ./...

.PHONY: check
check: ## run cyclic, security, performance, etc checks
	@printf "$(FMT_PRFX) running cyclic analysis\n"
	@gocyclo -over 16 -ignore ".cache/|vendor/|scripts/|docs/|deployment/|data/" .
	@printf "$(FMT_PRFX) running static security analysis\n"
	@gosec -tests -fmt=json -quiet -exclude-dir=vendor -exclude-dir=.cache -exclude-dir=scripts -exclude-dir=docs -exclude-dir=deployment ./...

.PHONY: align
align: ## align struct fields to use less memory
	@printf "$(FMT_PRFX) checking struct field memory alignment\n"
	@$(GO) list -f '{{.Dir}}' ./... | grep -v /vendor/ | \
		xargs fieldalignment ; if [[ $$? -eq 1 ]]; then  \
			printf "$(FMT_PRFX) $(FMT_WARN)unaligned struct fields detected$(FMT_END), check above output\n"; \
			printf "$(FMT_PRFX) to auto-fix run $(FMT_INFO)make align-fix$(FMT_END)\n"; \
		fi
		@printf "$(FMT_PRFX) $(FMT_OK)ok$(FMT_END)\n"; \

.PHONY: align-fix
align-fix: ## autofix misaligned struct fields
		@printf "$(FMT_PRFX) fixing struct field memory alignment\n"
		@$(GO) list -f '{{.Dir}}' ./... | grep -v /vendor/ | xargs fieldalignment -fix || exit 0;
		@printf "$(FMT_PRFX) aligned above files\n"
		@printf "$(FMT_PRFX) re-running $(FMT_INFO)make align$(FMT_END) to check for stragglers\n"
		@make align

.PHONY: fmt
fmt: ## format source files using gofumpt
	@printf "$(FMT_PRFX) formatting go files\n"
	@gofumpt -w $(SOURCE_FILES)

.PHONY: dev-deps-file
dev-deps-file: $(DIR_OUT) ## create a file with all dev deps listed
	@rm -f $(DIR_OUT)/dev_deps.txt
	@touch $(DIR_OUT)/dev_deps.txt
	@for tool in $(DEV_EXTERNAL_TOOLS) ; do \
  		printf "$$tool\n" >> $(DIR_OUT)/dev_deps.txt ; \
	done

.PHONY: dev-deps
dev-deps: ## install developer/ci dependencies
	@printf "$(FMT_PRFX) installing development/CI dependencies\n"
	@printf "$(FMT_PRFX) at $(FMT_INFO)$$(go env GOBIN)$(FMT_END)\n"
	@for tool in  $(DEV_EXTERNAL_TOOLS) ; do \
		printf "$(FMT_PRFX) installing/updating: $(FMT_INFO)$$tool$(FMT_END)\n" ; \
		$(GO) install $$tool; \
	done

CI_DEPENDENCIES_URL ?= https://github.com/sudoless/actions/releases/download/go-ci-deps/v0.1.0/
CI_DEPENDENCIES_TARGET ?= $$(go env GOOS)_$$(go env GOARCH).tar.gz
.PHONY: ci-deps
ci-deps: ## pull pre-packaged ci dependencies at GOBIN, using GOOS and GOARCH
	@printf "$(FMT_PRFX) pulling development/CI dependencies\n"
	@printf "$(FMT_PRFX) at $(FMT_INFO)$$(go env GOBIN)$(FMT_END)\n"
	@printf "$(FMT_PRFX) for $(FMT_INFO)$$(go env GOOS) $$(go env GOARCH)$(FMT_END)\n"
	@curl -sL -o $(CI_DEPENDENCIES_TARGET) $(CI_DEPENDENCIES_URL)$(CI_DEPENDENCIES_TARGET)
	@printf "$(FMT_PRFX) pulling checksum\n"
	@curl -sL -o $(CI_DEPENDENCIES_TARGET).sha256 $(CI_DEPENDENCIES_URL)$(CI_DEPENDENCIES_TARGET).sha256
	@$(SHA256) --check $(CI_DEPENDENCIES_TARGET).sha256
	@printf "$(FMT_PRFX) unpacking\n"
	@tar -xvf $(CI_DEPENDENCIES_TARGET) --directory=$$(go env GOBIN)
	@printf "$(FMT_PRFX) cleanup\n"
	@rm -rf $(CI_DEPENDENCIES_TARGET) $(CI_DEPENDENCIES_TARGET).sha256


# INTERNAL

.PHONY: mk-update
mk-update: ## update this Makefile, use THIS_MAKEFILE_UPDATE=... to specify version
	@printf "$(FMT_PRFX) updating this makefile from $(FMT_INFO)$(THIS_MAKEFILE_VERSION)$(FMT_END) to $(FMT_INFO)$(THIS_MAKEFILE_UPDATE)$(FMT_END)\n"
	@curl -s $(THIS_MAKEFILE_URL) > $(THIS_MAKEFILE).new
	@awk '/^#### CUSTOM/,0' $(THIS_MAKEFILE) | tail -n +2 >> $(THIS_MAKEFILE).new
	@mv -f $(THIS_MAKEFILE).new $(THIS_MAKEFILE)

.PHONY: help
help:
	@grep -h -E '^[a-zA-Z/_-]+%?:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m$(THIS_IMPORT_PREFIX)%-30s\033[0m %s\n", $$1, $$2}'


#### CUSTOM # Anything under the CUSTOM line is migrated by the mk-update command to the new Makefile version
