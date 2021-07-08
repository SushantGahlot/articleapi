#!/bin/bash
set -e

echo "#################### Migrate Database ####################"
go run cmd/migrate/main.go

echo "######### Download CompileDaemon for hot-reload ##########"

# Do not need to add it as dependency in go.mod
GO111MODULE=off go get github.com/githubnemo/CompileDaemon

echo "#################### Starting Daemon ####################"
CompileDaemon --build="go build -o main cmd/api/main.go" --command=./main
