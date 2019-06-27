GO111MODULE ?= on
GOOPTS := -mod=vendor

all: vendor build test

.PHONY: build
build:
	GO111MODULE=$(GO111MODULE) go build $(GOOPTS) ./cmd/remotewrite

.PHONY: test
test:
	GO111MODULE=$(GO111MODULE) go test $(GOOPTS) -race ./...

vendor: go.mod go.sum
	GO111MODULE=$(GO111MODULE) go mod vendor

.PHONY: docker-build
docker-build:
	docker build . -t csmarchbanks/remote-write-sidecar
