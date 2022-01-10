#      @ SUDOLESS SRL <contact@sudoless.org>
#      This Source Code Form is subject to the
#      terms of the Mozilla Public License, v.
#      2.0. If a copy of the MPL was not
#      distributed with this file, You can
#      obtain one at
#      http://mozilla.org/MPL/2.0/.


THIS_MAKEFILE_VERSION = v0.0.3
THIS_MAKEFILE_UPDATE = master
THIS_MAKEFILE := $(lastword $(MAKEFILE_LIST))
THIS_MAKEFILE_URL_BASE := https://raw.githubusercontent.com/sudoless/make/$(THIS_MAKEFILE_UPDATE)
THIS_MAKEFILE_URL := $(THIS_MAKEFILE_URL_BASE)/Makefile


# PATH
export PATH := $(abspath bin/):${PATH}

# META
export PROJECT_NAME := $(shell basename $(abspath .))
export PROJECT_ORG := $(shell basename $(abspath ..))

# META - FMT
export FMT_MISC := \033[90;1m
export FMT_INFO := \033[94;1m
export FMT_OK   := \033[92;1m
export FMT_WARN := \033[33;1m
export FMT_END  := \033[0m
export FMT_PRFX := $(FMT_MISC)=>$(FMT_END)


# GIT
ifneq ("$(wildcard .git/)","") # check .git/ exists
export GIT_TAG_HASH := $(shell git rev-list --abbrev-commit --tags --max-count=1)
export GIT_TAG := $(shell git describe --abbrev=0 --tags ${GIT_TAG_HASH} 2>/dev/null || true)
export GIT_VERSION := $(GIT_TAG)
export GIT_LATEST_HASH := $(shell git rev-parse --short HEAD)
export GIT_LATEST_COMMIT_DATE := $(shell git log -1 --format=%cd --date=format:"%Y%m%d")
export GIT_CHANGES := $(shell git rev-list $(GIT_TAG)..HEAD --count)

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
export SV_VERSION         	:= $(subst v,,$(GIT_TAG))
export SV_VERSION_PARTS   	:= $(subst ., ,$(SV_VERSION))
export SV_MAJOR           	:= $(word 1,$(SV_VERSION_PARTS))
export SV_MINOR           	:= $(word 2,$(SV_VERSION_PARTS))
export SV_MICRO           	:= $(word 3,$(SV_VERSION_PARTS))
export SV_MAJOR_NEXT      	:= $(shell echo $$(($(SV_MAJOR)+1)))
export SV_MINOR_NEXT      	:= $(shell echo $$(($(SV_MINOR)+1)))
export SV_MICRO_NEXT_1    	:= $(shell echo $$(($(SV_MICRO)+1)))
export SV_MICRO_NEXT      	:= $(shell echo $$(($(SV_MICRO)+$(GIT_CHANGES))))
export SV_GIT_MSG 			:= 'Bumping'
export SV_GIT_FLAGS			:= -a -m $(SV_GIT_MSG)

# BUILD
export BUILD_HASH		?= $(GIT_LATEST_HASH)
export BUILD_TIME		?= $$(date +%s)
export BUILD_VERSION	?= $(GIT_VERSION)


# IMPORTS
THIS_IMPORT_DIR ?= ./make
THIS_IMPORT_EXT ?= mk

IMPORTS ?= $(subst .$(THIS_IMPORT_EXT),,$(notdir $(wildcard $(THIS_IMPORT_DIR)/*.$(THIS_IMPORT_EXT))))

$(THIS_IMPORT_DIR):
	@mkdir -p $(THIS_IMPORT_DIR)

.PHONY: imports
imports: ## list existing imports, imports can be used as <IMPORT>/<rule>, eg: go/help
	@printf  "$(FMT_OK)$(IMPORTS)$(FMT_END)\n"

.PHONY: add/%
add/%: $(THIS_IMPORT_DIR)/%.$(THIS_IMPORT_EXT) ## add a new "import" makefile
	@printf "$(FMT_PRFX) done\n"

# not phony, because we do not want to overwrite existing imports
.PRECIOUS: $(THIS_IMPORT_DIR)/%.$(THIS_IMPORT_EXT)
$(THIS_IMPORT_DIR)/%.$(THIS_IMPORT_EXT): $(THIS_IMPORT_DIR)
	@printf "$(FMT_PRFX) adding import $(FMT_OK)$*$(FMT_END)\n"
	@printf "$(FMT_PRFX) downloading from $(FMT_INFO)$(THIS_MAKEFILE_URL_BASE)/$*.$(THIS_IMPORT_EXT)$(FMT_END)\n"
	@printf "$(FMT_PRFX) downloading to $(FMT_INFO)$(THIS_IMPORT_DIR)/$*.$(THIS_IMPORT_EXT)$(FMT_END)\n"
	@curl -s --fail-with-body $(THIS_MAKEFILE_URL_BASE)/$*.$(THIS_IMPORT_EXT) > $(THIS_IMPORT_DIR)/$*.$(THIS_IMPORT_EXT)


.PHONY: $(addsuffix /%,$(IMPORTS))
$(addsuffix /%,$(IMPORTS)):
	@THIS_IMPORT_PREFIX=$(firstword $(subst /, ,$@))/ $(MAKE) -f $(THIS_IMPORT_DIR)/$(firstword $(subst /, ,$@)).$(THIS_IMPORT_EXT) $*


.PHONY: info
info: ## display project information
	@printf "$(FMT_PRFX) printing info\n"
	@printf "$(FMT_PRFX) project name $(FMT_INFO)$(PROJECT_NAME)$(FMT_END)\n"
	@printf "$(FMT_PRFX) project org  $(FMT_INFO)$(PROJECT_ORG)$(FMT_END)\n"
	@printf "$(FMT_PRFX) project mod name $(FMT_INFO)$(PROJECT_MOD_NAME)$(FMT_END)\n"
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

.PHONY: info-version
info-version: ## prints the BUILD_VERSION and nothing else
	@printf "$(BUILD_VERSION)"

.PHONY: tag-micro tag-z
tag-z: tag-micro
tag-micro: ## tag the current commit with the next vX.Y.Z by adding number of git changes to Z
	@printf "$(FMT_PRFX) bumping $(FMT_INFO)$(GIT_TAG)$(FMT_END) to $(FMT_INFO)v$(SV_MAJOR).$(SV_MINOR).$(SV_MICRO_NEXT)$(FMT_END)\n"
	@git tag $(SV_GIT_FLAGS) v$(SV_MAJOR).$(SV_MINOR).$(SV_MICRO_NEXT)

.PHONY: tag-micro-one
tag-micro-one: ## tag the current commit with the next vX.Y.Z by adding 1 to Z
	@printf "$(FMT_PRFX) bumping $(FMT_INFO)$(GIT_TAG)$(FMT_END) to $(FMT_INFO)v$(SV_MAJOR).$(SV_MINOR).$(SV_MICRO_NEXT_1)$(FMT_END)\n"
	@git tag $(SV_GIT_FLAGS) v$(SV_MAJOR).$(SV_MINOR).$(SV_MICRO_NEXT_1)


.PHONY: tag-minor tag-y
tag-y: tag-minor
tag-minor: ## tag the current commit with the next vX.Y.Z by adding 1 to Y
	@printf "$(FMT_PRFX) bumping $(FMT_INFO)$(GIT_TAG)$(FMT_END) to $(FMT_INFO)v$(SV_MAJOR).$(SV_MINOR_NEXT).0$(FMT_END)\n"
	@git tag $(SV_GIT_FLAGS) v$(SV_MAJOR).$(SV_MINOR_NEXT).0

.PHONY: tag-major tag-x
tag-x: tag-major
tag-major: ## tag the current commit with the next vX.Y.Z by adding 1 to X
	@printf "$(FMT_PRFX) bumping $(FMT_INFO)$(GIT_TAG)$(FMT_END) to $(FMT_INFO)v$(SV_MAJOR_NEXT).0.0$(FMT_END)\n"
	@git tag $(SV_GIT_FLAGS) v$(SV_MAJOR_NEXT).0.0


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
