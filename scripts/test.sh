#!/bin/sh

set -e

echo "\033[36;1m*** Format check\033[0m"
gofmt -l .

echo "\033[36;1m*** Vet check\033[0m"
go vet ./...

echo "\033[36;1m*** Lint check\033[0m"
$HOME/go/bin/golint ./...

echo "\033[36;1m*** Tests\033[0m"
go test -v ./...
