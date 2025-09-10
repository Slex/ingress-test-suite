APP_NAME := ingress-test-suite
VERSION := $(if $(VERSION),$(VERSION),0.0.1)
BIN_DIR := bin
BIN_PATH := $(BIN_DIR)/$(APP_NAME)
DOCKER_IMAGE := $(APP_NAME):$(VERSION)

# Go настройки
GO := go
GO_FLAGS := -ldflags="-s -w -extldflags '-static' -X 'main.version=${VERSION}'"
GO_BUILD := CGO_ENABLED=0 $(GO) build $(GO_FLAGS) -o $(BIN_PATH)

.PHONY: build
build:
	@mkdir -p $(BIN_DIR)
	$(GO_BUILD)
	@echo "✅ Build done: $(BIN_PATH)"

.PHONY: clean
clean:
	@rm -rf $(BIN_DIR)
	@echo "🧹 Clean finished."

.PHONY: docker-build
docker-build: build
	@docker build -t $(DOCKER_IMAGE) .
	@echo "🐳 Docker build finish: $(DOCKER_IMAGE)"

.PHONY: all
all: clean build docker-build
	@echo "🎯 Rebuild all"
