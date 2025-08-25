#!/bin/bash

# Ask for package name if not passed as argument
if [ -z "$1" ]; then
  read -p "Enter the Go package to test (e.g. ./... or ./mypkg): " PKG
else
  PKG=$1
fi

# Coverage profile file
COVERFILE=coverprofile.out

# Run tests with coverage
go test "$PKG" -coverprofile=$COVERFILE

# Check if coverage file was created
if [ -f "$COVERFILE" ]; then
  # Show coverage in terminal
  go tool cover -func=$COVERFILE

  # Generate HTML report
  go tool cover -html=$COVERFILE -o coverage.html

  # Open in browser (Linux/macOS auto-detect)
  if command -v xdg-open >/dev/null; then
    xdg-open coverage.html >/dev/null 2>&1 &
  elif command -v open >/dev/null; then
    open coverage.html
  else
    echo "Coverage report generated: coverage.html"
    echo "Open it manually in your browser."
  fi
else
  echo "Coverage profile not generated."
fi
