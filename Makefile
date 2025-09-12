.PHONY: build run latest k8s

build:
	go build -ldflags="-s -w" -o dist/aegisbot .

run:
	$(MAKE) build
	./dist/aegisbot

latest:
	docker build -t kkrypt0nn/aegisbot:latest .
