#!/bin/bash

source .github/env.sh

go run ./infra/build/build.go
sudo cp -f build/current/sagernet /usr/local/bin/
