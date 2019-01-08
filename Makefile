.PHONY: all

CURDIR := $(shell pwd)

help: ## Print this help
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

serve: ## Start local server
	PORT=8081 go run ./backend/main.go &
	PORT=8082 TARGET_HOST=localhost:8081 TARGET_SCHEMA=http go run ./proxy/main.go &
	trap "kill %1 && kill %2" 1 2 3 15
	PORT=8080 BACKEND_HOST=localhost:8082 BACKEND_SCHEMA=http go run ./front/main.go
