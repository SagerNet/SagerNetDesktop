#!/bin/bash

source .github/env.sh

if [[ "$@" == "all" ]]; then

  go run ./infra/build/build.go all

else

  go run ./infra/build/build.go

fi
