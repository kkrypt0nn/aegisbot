.PHONY: build run latest proto fmt ensure-goimports

build:
	go build -ldflags="-s -w" -o dist/aegisbot ./cmd/aegisbot

run:
	$(MAKE) build
	./dist/aegisbot

latest:
	docker build -t kkrypt0nn/aegisbot:latest .

ensure-goimports:
	@command -v goimports >/dev/null 2>&1 || { \
		echo "Installing goimports..."; \
		go install golang.org/x/tools/cmd/goimports@latest; \
	}

fmt: ensure-goimports
	@echo "Formatting project..."
	@goimports -w .

PROTO_FILES := $(wildcard proto/*.proto)
proto:
	@for file in $(PROTO_FILES); \
		do protoc --go_out=. --go_opt=paths=source_relative $${file}; \
	done
	@echo "Formatting generated Go files..."
	@$(MAKE) fmt
