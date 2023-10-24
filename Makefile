export DATE := $(shell date +%Y.%m.%d-%H%M)
export LATEST_COMMIT := $(shell git log --pretty=format:'%h' -n 1)
export BRANCH := $(shell git branch |grep -v "no branch"| grep \*|cut -d ' ' -f2)
export GIT_REPO := $(shell git config --get remote.origin.url  2> /dev/null)
export BIN_DIR=./bin


export VERSION_FILE   := version.txt
export TAG     := $(shell [ -f "$(VERSION_FILE)" ] && cat "$(VERSION_FILE)" || echo '0.5.46')
export VERMAJMIN      := $(subst ., ,$(TAG))
export VERSION        := $(word 1,$(VERMAJMIN))
export MAJOR          := $(word 2,$(VERMAJMIN))
export MINOR          := $(word 3,$(VERMAJMIN))
export NEW_MINOR      := $(shell expr "$(MINOR)" + 1)
export NEW_TAG := $(VERSION).$(MAJOR).$(NEW_MINOR)


export PROTOC2=protoc-2.6.1
export PROTOC3=protoc

ifeq ($(BRANCH),)
BRANCH := master
endif

export COMPILE_LDFLAGS=-s -X "main.BuildDate=${DATE}" \
                          -X "main.LatestCommit=${LATEST_COMMIT}" \
						  -X "main.Version=${NEW_TAG}"\
						  -X "main.GitRepo=${GIT_REPO}" \
                          -X "main.GitBranch=${BRANCH}"



build_info:  ## Build the container
	@echo ''
	@echo '---------------------------------------------------------'
	@echo 'DATE              $(DATE)'
	@echo 'LATEST_COMMIT     $(LATEST_COMMIT)'
	@echo 'BRANCH            $(BRANCH)'
	@echo 'COMPILE_LDFLAGS   $(COMPILE_LDFLAGS)'
	@echo 'TAG              $(TAG)'
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

protogen: build_info ## build protogen binary in bin dir
	make BIN_NAME=protogen APP_PATH=github.com/alexj212/protogen build_app
	@echo ''
	@echo 'You can now copy the binary `protogen` into the system path'
	@echo ''

install: build_info ## install protogen binary in $GOPATH/bin dir
	go install github.com/alexj212/protogen

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




####################################################################################################################
##
## Code vetting tools
##
####################################################################################################################


gotest: ## run tests
	go test -v $(PROJ_PATH)/...

fmt: ## run fmt on project
	#go fmt $(PROJ_PATH)/...
	gofmt -s -d -w -l .

doc: ## launch godoc on port 6060
	godoc -http=:6060

deps: ## display deps for project
	go list -f '{{ join .Deps  "\n"}}' . |grep "/" | grep -v $(PROJ_PATH)| grep "\." | sort |uniq

lint: ## run lint on the project
	golint ./...

staticcheck: ## run staticcheck on the project
	staticcheck -ignore "$(shell cat .checkignore)" .

vet: ## run go vet on the project
	go vet .

reportcard: fmt ## run goreportcard-cli
	goreportcard-cli -v

tools: ## install dependent tools for code analysis
	go install github.com/gordonklaus/ineffassign@latest
	go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
	go install golang.org/x/lint/golint@latest
	go install github.com/gojp/goreportcard/cmd/goreportcard-cli@latest
	go install github.com/goreleaser/goreleaser@latest


upgrade:
	go get -u ./...
	go mod tidy







publish: ## tag & push to gitlab
	@echo "\n\n\n\n\n\nRunning git add\n"
	echo "$(NEW_TAG)" > "$(VERSION_FILE)"
	git add -A
	@echo "\n\n\n\n\n\nRunning git commit v$(NEW_TAG)\n"
	git commit -m "latest version: v$(NEW_TAG)"

	@echo "\n\n\n\n\n\nRunning git tag\n"
	git tag  "v$(NEW_TAG)"

	@echo "\n\n\n\n\n\nRunning git push\n"
	git push -f origin "v$(NEW_TAG)"

	git push -f



