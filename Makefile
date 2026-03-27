.PHONY: all build frontend clean install-deps

BINARY_NAME := bedrock-timeline
FRONTEND_DIR := web
BUILD_DIR := build

all: build

install-deps:
	go mod download
	go mod tidy
	cd $(FRONTEND_DIR) && npm install

frontend:
	cd $(FRONTEND_DIR) && npm run build
	mkdir -p cmd/server/static
	cp -r $(FRONTEND_DIR)/static/* cmd/server/static/

build: frontend
	CGO_ENABLED=1 go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/server

clean:
	rm -rf $(BUILD_DIR)
	rm -rf cmd/server/static
	rm -rf $(FRONTEND_DIR)/.svelte-kit
	rm -rf $(FRONTEND_DIR)/static

dev:
	cd $(FRONTEND_DIR) && npm run dev

run:
	go run ./cmd/server

install:
	install -m 755 $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	install -m 644 bedrock-timeline.service /etc/systemd/system/
	mkdir -p /opt/bedrock-timeline/data
	systemctl daemon-reload