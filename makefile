SHELL := /bin/bash

run:
	go run .

tidy:
	go mod tidy

docker:
	docker buildx build --platform=linux/amd64,linux/arm64 -t wonko/ingress-whitelist:0.0.2 -t wonko/ingress-whitelist:latest --push .
