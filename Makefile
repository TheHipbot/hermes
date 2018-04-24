# Environment info
GOCMD := go
SOURCEDIR=.
SOURCES := $(find $(SOURCEDIR) -name '*.go')

# Binary info
BINARY := hermes
BIN_DIR=
PKG_NAME := hermes
PKG_PREFIX := "github.com/TheHipbot/${PKG_NAME}"

GOOS := darwin
GOARCH := amd64

# Build info
VERSION := 0.0.1
BUILD := `git rev-parse --short HEAD`

# OS's and Architectures
goos := darwin \
		linux \
		windows

archs := 386 \
		 amd64

.DEFAULT_GOAL: $(BINARY)

.PHONY: clean
clean:
	$(GOCMD) clean && if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi && rm -rf bin

.PHONY: lint
lint: 
	golint ./...

.PHONY: vet
vet:
	$(GOCMD) vet ./...

.PHONY: install
install: ensure
	$(GOCMD) install

.PHONY: ensure
ensure:
	dep ensure

build: ensure $(SOURCES)
	GOOS=${GOOS} GOARCH=${GOARCH} $(GOCMD) build -o ${BIN_DIR}${BINARY}

.PHONY: test
test: 
	go test ./...

cross-compile: all-archs

all-archs: $(archs)

$(archs):
	$(MAKE) GOARCH=$@ all-os

all-os: $(goos)

$(goos):
	$(MAKE) GOOS=$@ BIN_DIR=bin/ BINARY=${BINARY}-$@-${GOARCH} build


.PHONY: docker-build
docker-build:
	docker build -t "thehipbot/${PKG_NAME}:${BUILD}" --build-arg PACKAGE=$(PKG_PREFIX) .