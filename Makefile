.PHONY: all
all: help

mod:	## Go mod
	@echo "go mod..."
	@go mod tidy
	@go mod verify

help:	## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


