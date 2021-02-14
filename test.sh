#!/bin/bash

# Generate Mocks
bash generate_mocks.sh

# Run test
go test -v --coverprofile=coverage.out --coverpkg ./infra,./service,./repository,./controller,./contract/...,./middleware,./router,./infra,./dto,./model ./test/unit ./test/integration

# Get coverage html
go tool cover -html=coverage.out -o cover.html