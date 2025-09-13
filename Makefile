.PHONY: build run latest proto

build:
	go build -ldflags="-s -w" -o dist/aegisbot .

run:
	$(MAKE) build
	./dist/aegisbot

latest:
	docker build -t kkrypt0nn/aegisbot:latest .

PROTO_FILES := $(wildcard proto/*.proto)
proto:
	@for file in $(PROTO_FILES); \
		do protoc --go_out=. --go_opt=paths=source_relative $${file}; \
	done
