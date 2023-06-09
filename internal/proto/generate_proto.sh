#!/usr/bin/env sh

protoc --go_out=./internal/generated/ --go-grpc_out=./internal/generated/ ./internal/proto/*.proto
