# Variables
TRG_PKG := app
BUILD_TIME := $(shell date +"%Y%m%d.%H%M%S")
GoVersion := N/A
CommitHash := N/A

# Determine Go version, Git tag, and Commit hash
GoVersion := $(shell go version | cut -d' ' -f3 | tr -d 'go')
CommitHash := $(shell git log -1 --pretty=format:%h 2>/dev/null || echo 'N/A')

# Flags for Go build/install
FLAG = -ldflags="-X 'main.BuildTime=$(BUILD_TIME)' \
                -X 'main.CommitHash=$(CommitHash)' \
                -X 'main.GoVersion=$(GoVersion)'"

# Default target
.PHONY: build

build:
	@echo 'go build'
	GOOS=linux GOARCH=arm64 go build -v $(FLAG) -o $(TRG_PKG)

clean:
	@echo 'Cleaning up...'
	rm -f ./$(TRG_PKG)
