@echo off
set DATABASE_URL=postgresql://postgres:123@localhost:5432/ParkingManagementTest
go run ./cmd/web/main.go
