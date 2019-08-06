export DATE := $(shell date +%Y.%m.%d-%H%M)
export LATEST_COMMIT := $(shell git log --pretty=format:'%h' -n 1)
export BRANCH := $(shell git branch |grep -v "no branch"| grep \*|cut -d ' ' -f2)
export BUILT_ON_IP := $(shell [ $$(uname) = Linux ] && hostname -i || hostname )

export BUILT_ON_OS=$(shell uname -a)
ifeq ($(BRANCH),)
BRANCH := master
endif

export COMMIT_CNT := $(shell git rev-list HEAD | wc -l | sed 's/ //g' )
export BUILD_NUMBER := ${BRANCH}-${COMMIT_CNT}
export COMPILE_LDFLAGS=-s -X "main.DATE=${DATE}" \
                          -X "main.LATEST_COMMIT=${LATEST_COMMIT}" \
                          -X "main.BUILD_NUMBER=${BUILD_NUMBER}" \
                          -X "main.BUILT_ON_IP=${BUILT_ON_IP}" \
                          -X "main.BUILT_ON_OS=${BUILT_ON_OS}"

build_info: ## Build the container
	@echo ''
	@echo '---------------------------------------------------------'
	@echo 'BUILT_ON_IP      $(BUILT_ON_IP)'
	@echo 'BUILT_ON_OS      $(BUILT_ON_OS)'
	@echo 'DATE             $(DATE)'
	@echo 'LATEST_COMMIT    $(LATEST_COMMIT)'
	@echo 'BRANCH           $(BRANCH)'
	@echo 'COMMIT_CNT       $(COMMIT_CNT)'
	@echo 'BUILD_NUMBER     $(BUILD_NUMBER)'
	@echo 'COMPILE_LDFLAGS  $(COMPILE_LDFLAGS)'
	@echo 'PATH             $(PATH)'
	@echo '---------------------------------------------------------'
	@echo ''


####################################################################################################################
##
## help for each task - https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
##
####################################################################################################################
.PHONY: help

help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help



####################################################################################################################
##
## Build of binaries
##
####################################################################################################################

binaries: protogen ## build binaries in bin dir

build_app:
	go build -o ./$(BIN_NAME) -a -ldflags '$(COMPILE_LDFLAGS)' $(APP_PATH)

protogen: build_info ## build broadcastclient binary in bin dir
	make BIN_NAME=protogen APP_PATH=github.com/alexj212/protogen build_app
	@echo ''
	@echo 'You can now copy the binary `protogen` into the system path'
	@echo ''



####################################################################################################################
##
## Cleanup of binaries
##
####################################################################################################################

clean_binaries: protogen  ## clean all binaries in bin dir

clean_protogen: ## clean protogen
	rm -f ./protogen



test: ## clean protogen
	protoc  --go_out=./   ./_test/test.proto
	./protogen ./_test/test.proto Packet ./_test/mapping.go
