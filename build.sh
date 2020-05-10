#!/usr/bin/env bash

GOOS=linux GOARCH=arm GOARM=5 go build -ldflags "-s -w"