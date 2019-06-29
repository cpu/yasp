#!/bin/bash

set -e

export GO111MODULE=on

golangci-lint run && \
  go test -v ./... && \
    go install -race ./... && \
     GORACE="halt_on_error=1" yasp
