#!/bin/bash
project="kohme"

go mod tidy
go build -ldflags "-s -w" -o $project ./cmd/bot


docker build -t $project .
rm -rf $project
echo "build success $project:$version"