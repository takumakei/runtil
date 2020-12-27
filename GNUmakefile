.DEFAULT_GOAL := help

.PHONY: help build archive

SNAPSHOT := $(shell [[ -z $$(git status --porcelain) ]] || echo --snapshot)

help: ## show this message
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-8s\033[0m %s\n", $$1, $$2}'

build: ## build
	goreleaser build --rm-dist $(SNAPSHOT)

archive: ## archive
	goreleaser --rm-dist --skip-publish $(SNAPSHOT)
