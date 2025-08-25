#!/bin/bash

export DATABASE_URL="postgresql://kaushik:123@localhost:5432/ParkingManagement"

go run ./cmd/web/main.go
