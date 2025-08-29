#!/bin/bash

#| sed 's/,$//'

go test -cover -coverpkg=$(go list ./... | grep -v -E "/mocks" | tr '\n' ',' ) -coverprofile=cover2.out ./...
go tool cover -func=cover2.out