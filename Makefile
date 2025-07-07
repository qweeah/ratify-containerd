WATCHER_BINARY_NAME = configmap-watcher
VERIFIER_BINARY_NAME = verifier

all: build

.PHONY: build
build: build-watcher build-verifier

.PHONY: build-watcher
build-watcher: fmt vet
	go build -o ./bin/${WATCHER_BINARY_NAME} ./cmd/watcher

.PHONY: build-verifier
build-verifier: fmt vet
	go build -o ./bin/${VERIFIER_BINARY_NAME} ./cmd/verifier

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...