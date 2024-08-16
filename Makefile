PACKAGE := github.com/marcsello/ponyhug2-backend

# define the build timestamp, commit hash, version
BUILD_TIMESTAMP := $(shell date '+%Y-%m-%dT%H:%M:%S')
COMMIT_HASH ?= $(DRONE_COMMIT)
BUILD_VERSION ?= $(DRONE_BUILD_NUMBER)

ifeq ($(BUILD_VERSION),)
	# fallback to dev when not built with drone
	BUILD_VERSION := dev
endif

ifeq ($(COMMIT_HASH),)
	# fallback to git to figure this out (might not work if this is not a git repo)
	COMMIT_HASH := $(shell git rev-parse --short HEAD)
endif

# define build parameters
export GOOS := linux
LDFLAGS := -X 'main.version=${BUILD_VERSION}' -X 'main.commitHash=${COMMIT_HASH}' -X 'main.buildTimestamp=${BUILD_TIMESTAMP}'

.PHONY: all
all: ponyhug

ponyhug: main.go dist
	GOARCH=amd64 go build -v -ldflags="${LDFLAGS}" -o "dist/ponyhug" "."

dist:
	mkdir -p dist

.PHONY: clean
clean:
	rm -rf dist/
