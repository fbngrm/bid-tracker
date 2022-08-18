#!/usr/bin/env just --justfile

default:
	just -l

build: generate
	cd server && go build
	@echo "server built"
	cd cmd/client && go build
	@echo "client built"
	@just lint

lint:
	golangci-lint run

release:
	just buf breaking --against ".git#branch=main,subdir=."

setup: setup-go generate

generate:
	make gen/proto/go/auction/v1/*.go >/dev/null

run: build
    docker-compose up

run-server:
	cd server
	./server/server

client: generate
	cd cmd/client && go build
	./cmd/client/client

buf *args='':
	cd apis && PATH="$PATH:$(go env GOPATH)/bin" buf {{args}}

setup-go:
	cd gen && go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	cd gen && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	cd gen && go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
	cd gen && go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
