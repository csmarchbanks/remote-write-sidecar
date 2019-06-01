.PHONY := build test docker

GO111MODULE ?= on
GOOPTS := -mod=vendor

all: build test

build:
	GO111MODULE=$(GO111MODULE) go build $(GOOPTS) ./cmd/remotewrite

test:
	GO111MODULE=$(GO111MODULE) go test $(GOOPTS) -race ./...

docker:
	docker build . -t csmarchbanks/remote-write-sidecar
