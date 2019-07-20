GO111MODULE ?= on
GOOPTS := -mod=vendor
DOCKER_REPO ?= quay.io/csmarchbanks/remote-write-sidecar
DOCKER_TAG ?= $(shell git rev-parse --abbrev-ref HEAD)-$(shell date -u +"%Y-%m-%d")-$(shell git rev-parse --short HEAD)

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
	docker build . -t $(DOCKER_REPO):$(DOCKER_TAG)

.PHONY: docker-publish
docker-publish: docker-build
	docker push $(DOCKER_REPO):$(DOCKER_TAG)
