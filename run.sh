#!/bin/bash

source .github/env.sh

go run ./infra/build/build.go
build/current/sagernet