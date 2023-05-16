VERSION ?= $(shell runner/gnmic version | cut -d ' ' -f 3 | cut -d ':' -f 1)
BUILD_DIR = build
CGO_ENABLED ?= 0
DOCKERHUB_REPO = ghcr.io/orb-community
GOARCH ?= $(shell dpkg-architecture -q DEB_BUILD_ARCH)
COMMIT_HASH = $(shell git rev-parse --short HEAD)
REF_TAG ?= latest

getgnmic:
	wget -O /tmp/gnmic.tar.gz https://github.com/openconfig/gnmic/releases/download/v0.30.0/gnmic_0.30.0_Linux_x86_64.tar.gz
	tar -xvzf /tmp/gnmic.tar.gz -C /tmp/
	mv /tmp/gnmic runner/gnmic
	rm -rf /tmp/gnmic*

build:
	CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=$(GOARCH) GOARM=$(GOARM) go build -o ${BUILD_DIR}/orb-gnmic cmd/main.go
  
container:
	docker build --no-cache \
	  --tag=$(DOCKERHUB_REPO)/orb-gnmic:$(REF_TAG) \
	  --tag=$(DOCKERHUB_REPO)/orb-gnmic:$(VERSION) \
	  --tag=$(DOCKERHUB_REPO)/orb-gnmic:$(VERSION)-$(COMMIT_HASH) \
	  -f docker/Dockerfile .
