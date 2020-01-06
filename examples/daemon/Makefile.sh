#!/usr/bin/env bash

rm -rf daemon-hook
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o daemon-hook main.go

docker build -t php-base .