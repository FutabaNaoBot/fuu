#!/bin/bash
project="kohme"

go mod tidy
go build -ldflags -s -ldflags -w -o $project .

docker build -t $project .
rm -rf $project
echo "build success $project:$version"