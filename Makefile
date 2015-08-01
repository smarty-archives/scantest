#!/usr/bin/make -f

# This Makefile is an example of what you could feed to scantest's -command flag.

default: test

test: build
	go test

cover: build
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out

build:
	go build
	go generate
