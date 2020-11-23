#!/bin/bash

richgo test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out