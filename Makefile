#      @ SUDOLESS SRL <contact@sudoless.org>
#      This Source Code Form is subject to the
#      terms of the Mozilla Public License, v.
#      2.0. If a copy of the MPL was not
#      distributed with this file, You can
#      obtain one at
#      http://mozilla.org/MPL/2.0/.

THIS_MAKEFILE_VERSION = v0.1.8
THIS_MAKEFILE_UPDATE = master
THIS_MAKEFILE := $(lastword $(MAKEFILE_LIST))
THIS_MAKEFILE_URL := https://raw.githubusercontent.com/sudoless/make/$(THIS_MAKEFILE_UPDATE)/golang.mk


# PATH
export PATH := $(abspath bin/):${PATH}

# META
ifneq ("$(wildcard go.mod/)","") # check go.mod exists
PROJECT_MOD_NAME := $(shell go list -m -mod=readonly)
PROJECT_NAME := $(notdir $(PROJECT_MOD_NAME))
endif

# META - FMT
FMT_MISC := \033[90;1m
FMT_INFO := \033[94;1m
FMT_OK   := \033[92;1m
FMT_WARN := \033[33;1m
FMT_END  := \033[0m
FMT_PRFX := $(FMT_MISC)=>$(FMT_END)

# GO
export CGO_ENABLED ?= 0
GO ?= GO111MODULE=on go
GO_TAGS ?= timetzdata

# OUTPUT
DIR_OUT   := out
FILE_COV  := $(DIR_OUT)/cover.out

# GIT
ifneq ("$(wildcard .git/)","") # check .git/ exists
GIT_TAG_HASH := $(shell git rev-list --abbrev-commit --tags --max-count=1)
GIT_TAG := $(shell git describe --abbrev=0 --tags ${GIT_TAG_HASH} 2>/dev/null || true)
GIT_VERSION := $(GIT_TAG)
GIT_LATEST_HASH := $(shell git rev-parse --short HEAD)
GIT_LATEST_COMMIT_DATE := $(shell git log -1 --format=%cd --date=format:"%Y%m%d")
GIT_CHANGES := $(shell git rev-list $(GIT_TAG)..HEAD --count)

ifneq ($(GIT_LATEST_HASH),$(GIT_TAG_HASH))
	GIT_VERSION := $(GIT_VERSION)-wip$(GIT_CHANGES).$(GIT_LATEST_HASH)
endif
ifeq ($(GIT_VERSION),)
	GIT_VERSION := -new.$(GIT_LATEST_HASH).$(GIT_LATEST_COMMIT_DATE)
endif
ifneq ($(shell git status --porcelain),)
	GIT_VERSION := $(GIT_VERSION)-dirty.$(GIT_LATEST_COMMIT_DATE).$$(whoami)
endif
endif

# SEMVER
SV_VERSION            := $(subst v,,$(GIT_TAG))
SV_VERSION_PARTS      := $(subst ., ,$(SV_VERSION))
SV_MAJOR              := $(word 1,$(SV_VERSION_PARTS))
SV_MINOR              := $(word 2,$(SV_VERSION_PARTS))
SV_MICRO              := $(word 3,$(SV_VERSION_PARTS))
SV_MAJOR_NEXT         := $(shell echo $$(($(SV_MAJOR)+1)))
SV_MINOR_NEXT         := $(shell echo $$(($(SV_MINOR)+1)))
SV_MICRO_NEXT_1       := $(shell echo $$(($(SV_MICRO)+1)))
SV_MICRO_NEXT         := $(shell echo $$(($(SV_MICRO)+$(GIT_CHANGES))))
SV_GIT_MSG := 'Bumping'
SV_GIT_FLAGS := -a -m $(SV_GIT_MSG)

# BUILD
BUILD_HASH ?= $(GIT_LATEST_HASH)
BUILD_TIME ?= $$(date +%s)
BUILD_VERSION ?= $(GIT_VERSION)

# SOURCE
SOURCE_FILES?=$$(find . -name '*.go' | grep -v pb.go | grep -v vendor)

# DEV - EXTERNAL TOOLS
DEV_EXTERNAL_TOOLS=\
	github.com/golangci/golangci-lint/cmd/golangci-lint@v1.39.0 \
	github.com/securego/gosec/v2/cmd/gosec@v2.7.0 \
	github.com/client9/misspell/cmd/misspell@v0.3.4 \
	github.com/fzipp/gocyclo/cmd/gocyclo@v0.3.1 \
	github.com/jstemmer/go-junit-report@v0.9.1 \
	golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@v0.1.0 \
	mvdan.cc/gofumpt@v0.1.1 \
	gotest.tools/gotestsum@v1.6.4

# DOCKER
THIS_DOCKER_BUILD_FLAGS ?=
THIS_DOCKER_DIR ?= ./deployment/docker
THIS_DOCKER_FILE ?= $(THIS_DOCKER_DIR)/Dockerfile
THIS_DOCKER_TAG ?= $(BUILD_VERSION)
THIS_DOCKER_IMAGE ?= $(PROJECT_MOD_NAME)/$*
THIS_DOCKER_ARTIFACT ?= $(THIS_DOCKER_IMAGE):$(THIS_DOCKER_TAG)


all: clean align spelling check lint test


.PHONY: info
info: ## display project information
	@printf "$(FMT_PRFX) printing info\n"
	@printf "$(FMT_PRFX) project mod name $(FMT_INFO)$(PROJECT_MOD_NAME)$(FMT_END)\n"
	@printf "$(FMT_PRFX) project name $(FMT_INFO)$(PROJECT_NAME)$(FMT_END)\n"
	@printf "$(FMT_PRFX) makefile $(FMT_INFO)$(THIS_MAKEFILE)$(FMT_END)\n"
	@printf "$(FMT_PRFX) makefile version $(FMT_INFO)$(THIS_MAKEFILE_VERSION)$(FMT_END)\n"
	@printf "$(FMT_PRFX) build hash $(FMT_INFO)$(BUILD_HASH)$(FMT_END)\n"
	@printf "$(FMT_PRFX) build current version $(FMT_INFO)$(BUILD_VERSION)$(FMT_END)\n"
	@printf "$(FMT_PRFX) git tag commit $(FMT_INFO)$(GIT_TAG_HASH)$(FMT_END)\n"
	@printf "$(FMT_PRFX) git tag $(FMT_INFO)$(GIT_TAG)$(FMT_END)\n"
	@printf "$(FMT_PRFX) git version $(FMT_INFO)$(GIT_VERSION)$(FMT_END)\n"
	@printf "$(FMT_PRFX) git latest commit $(FMT_INFO)$(GIT_LATEST_HASH)$(FMT_END)\n"
	@printf "$(FMT_PRFX) git latest commit date $(FMT_INFO)$(GIT_LATEST_COMMIT_DATE)$(FMT_END)\n"
	@printf "$(FMT_PRFX) git commit changes $(FMT_INFO)$(GIT_CHANGES)$(FMT_END)\n"


.PHONY: tag-micro
tag-micro: ## tag the current commit with the next vX.Y.Z by adding number of git changes to Z
	@printf "$(FMT_PRFX) bumping $(FMT_INFO)$(GIT_TAG)$(FMT_END) to $(FMT_INFO)v$(SV_MAJOR).$(SV_MINOR).$(SV_MICRO_NEXT)$(FMT_END)\n"
	@git tag $(SV_GIT_FLAGS) v$(SV_MAJOR).$(SV_MINOR).$(SV_MICRO_NEXT)

.PHONY: tag-micro-one
tag-micro-one: ## tag the current commit with the next vX.Y.Z by adding 1 to Z
	@printf "$(FMT_PRFX) bumping $(FMT_INFO)$(GIT_TAG)$(FMT_END) to $(FMT_INFO)v$(SV_MAJOR).$(SV_MINOR).$(SV_MICRO_NEXT_1)$(FMT_END)\n"
	@git tag $(SV_GIT_FLAGS) v$(SV_MAJOR).$(SV_MINOR).$(SV_MICRO_NEXT_1)


.PHONY: tag-minor
tag-minor: ## tag the current commit with the next vX.Y.Z by adding 1 to Y
	@printf "$(FMT_PRFX) bumping $(FMT_INFO)$(GIT_TAG)$(FMT_END) to $(FMT_INFO)v$(SV_MAJOR).$(SV_MINOR_NEXT).0$(FMT_END)\n"
	@git tag $(SV_GIT_FLAGS) v$(SV_MAJOR).$(SV_MINOR_NEXT).0

.PHONY: tag-major
tag-major: ## tag the current commit with the next vX.Y.Z by adding 1 to X
	@printf "$(FMT_PRFX) bumping $(FMT_INFO)$(GIT_TAG)$(FMT_END) to $(FMT_INFO)v$(SV_MAJOR_NEXT).0.0$(FMT_END)\n"
	@git tag $(SV_GIT_FLAGS) v$(SV_MAJOR_NEXT).0.0


.PHONY: init
init: ## setup a barebones Go project
	@mkdir -p cmd/ pkg/ docs/ scripts/ data/ deployment/ internal/
	@touch deployment/.netrc
	@echo "deployment/.netrc\n**/.env.*" > .gitignore

.PHONY: info-version
info-version: ## prints the BUILD_VERSION and nothing else
	@printf "$(BUILD_VERSION)"

.PHONY: run-%
run-%: build-% ## run the specified target
	@printf "$(FMT_PRFX) running $(FMT_INFO)$*$(FMT_END) from $(FMT_INFO)$(DIR_OUT)/dist/$*_$$(go env GOOS)_$$(go env GOARCH)$(FMT_END)\n"
	@$(DIR_OUT)/dist/$*_$$(go env GOOS)_$$(go env GOARCH)

.PHONY: build-%
build-%: APP_OUT ?= $*_$$(go env GOOS)_$$(go env GOARCH)
build-%: ## build a specific cmd/$(TARGET)/... into $(DIR_OUT)/dist/$(TARGET)...
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

.PHONY: install-%
install-%: APP_OUT ?= $*_$$(go env GOOS)_$$(go env GOARCH)
install-%: GOBIN ?= $(GOPATH)/bin
install-%: build-% ## install the built binary to $GOBIN
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

$(DIR_OUT):
	@mkdir -p $(DIR_OUT)

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
	@golangci-lint run -v --timeout 10m --skip-dirs=".cache/|vendor/|scripts/|docs/|deployment/"  ./...

.PHONY: check
check: ## run cyclic, security, performance, etc checks
	@printf "$(FMT_PRFX) running cyclic analysis\n"
	@gocyclo -over 16 -ignore ".cache/|vendor/|scripts/|docs/|deployment/" .
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

.PHONY: spelling
spelling: ## run misspell check
	@printf "$(FMT_PRFX) checking for spelling errors\n"
	@misspell -error pkg/
	@misspell -error cmd/

.PHONY: dev-deps
dev-deps: ## pull developer/ci dependencies
	@printf "$(FMT_PRFX) pulling development/CI dependencies\n"
	@for tool in  $(DEV_EXTERNAL_TOOLS) ; do \
		printf "$(FMT_PRFX) installing/updating: $(FMT_INFO)$$tool$(FMT_END)\n" ; \
		$(GO) install $$tool; \
	done

.PHONY: docker-list
docker-list: ## list docker images for the current project
	@printf "$(FMT_PRF) listing images for $(FMT_INFO)$(PROJECT_NAME)$(FMT_END) project\n"
	@docker images -f label=project=$(PROJECT_NAME)

.PHONY: docker-build-%
docker-build-%: ## build docker image
	@printf "$(FMT_PRFX) building with docker $(FMT_INFO)$$(docker version -f 'server: {{.Server.Version}}, client: {{.Client.Version}}')$(FMT_END)\n"
	@printf "$(FMT_PRFX) docker on host $(FMT_WARN)$(DOCKER_HOST)$(FMT_END)\n"
	@printf "$(FMT_PRFX) docker file $(FMT_INFO)$(THIS_DOCKER_FILE)$(FMT_END)\n"
	@printf "$(FMT_PRFX) docker artifact output $(FMT_INFO)$(THIS_DOCKER_ARTIFACT)$(FMT_END)\n"
	@DOCKER_BUILDKIT=1 docker build $(THIS_DOCKER_BUILD_FLAGS) \
		--secret id=netrc,src=./deployment/.netrc \
		--build-arg APP_NAME=$* \
		--build-arg BUILD_VERSION=$(BUILD_VERSION) \
		--build-arg BUILD_HASH=$(BUILD_HASH) \
		-f $(THIS_DOCKER_FILE) -t $(THIS_DOCKER_ARTIFACT) \
		--label "project=$(PROJECT_NAME)" \
		--label "build_hash=$(BUILD_HASH)" \
		--label "build_time=$(BUILD_TIME)" \
		--label "build_machine=$$(whoami)@$$(hostname)" .
	@printf "$(FMT_PRFX) docker artifact output $(FMT_INFO)$(THIS_DOCKER_ARTIFACT)$(FMT_END)\n"
	@printf "$(FMT_PRFX) run $(FMT_INFO)docker tag $(THIS_DOCKER_ARTIFACT) ...$(FMT_END) to change name\n"

.PHONY: docker-tag-%
docker-tag-%: ## tags the last built docker image for the given package using its version and $IMAGE_BASE
	@printf "$(FMT_PRFX) tagging $(FMT_INFO)$(THIS_DOCKER_ARTIFACT)$(FMT_END)\n"
	@printf "$(FMT_PRFX) as      $(FMT_INFO)$(IMAGE_BASE)/$*:$(THIS_DOCKER_TAG)$(FMT_END)\n"
	@docker tag $(THIS_DOCKER_ARTIFACT) $(IMAGE_BASE)/$*:$(THIS_DOCKER_TAG)

.PHONY: docker-push-%
docker-push-%: ## pushes the last tagged docker image for the given package using its version and $IMAGE_BASE
	@printf "$(FMT_PRFX) pushing $(FMT_INFO)$(IMAGE_BASE)/$*:$(THIS_DOCKER_TAG)$(FMT_END)\n"
	@docker push $(IMAGE_BASE)/$*:$(THIS_DOCKER_TAG)

.PHONY: mk-update
mk-update: ## update this Makefile, use THIS_MAKEFILE_UPDATE=... to specify version
	@printf "$(FMT_PRFX) updating this makefile from $(FMT_INFO)$(THIS_MAKEFILE_VERSION)$(FMT_END) to $(FMT_INFO)$(THIS_MAKEFILE_UPDATE)$(FMT_END)\n"
	@curl -s $(THIS_MAKEFILE_URL) > $(THIS_MAKEFILE).new
	@awk '/^#### CUSTOM/,0' Makefile | tail -n +2 >> $(THIS_MAKEFILE).new
	@mv -f Makefile.new Makefile

.PHONY: help
help:
	@grep -h -E '^[a-zA-Z_-]+%?:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


#### CUSTOM # Anything under the CUSTOM line is migrated by the mk-update command to the new Makefile version
