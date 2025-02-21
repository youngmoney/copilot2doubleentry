#!/usr/bin/env bash

GOOS=darwin GOARCH=arm64 go build -o bin/copilot2doubleentry-darwin-arm64 .
GOOS=linux GOARCH=amd64 go build -o bin/copilot2doubleentry-linux-amd64 .
