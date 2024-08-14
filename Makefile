PACKAGE := github.com/marcsello/ponyhug2-backend

# define the build timestamp, commit hash, version
BUILD_TIMESTAMP := $(shell date '+%Y-%m-%dT%H:%M:%S')
COMMIT_HASH := $(shell git rev-parse --short HEAD)
BUILD_VERSION ?= $(DRONE_TAG)

ifeq ($(BUILD_VERSION),)
	# fallback to the branch name
	BUILD_VERSION := $(shell git rev-parse --abbrev-ref HEAD)
endif

# define build parameters
export GOOS := linux
LDFLAGS := -X 'main.version=${BUILD_VERSION}' -X 'main.commitHash=${COMMIT_HASH}' -X 'main.buildTimestamp=${BUILD_TIMESTAMP}'

.PHONY: all
all: ponyhug

ponyhug: main.go dist
	GOARCH=amd64 go build -v -ldflags="${LDFLAGS}" -o "ponyhug" "."

dist:
	mkdir -p dist

.PHONY: clean
clean:
	rm -rf dist/