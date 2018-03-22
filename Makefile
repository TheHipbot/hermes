# Environment info
GOCMD := go
SOURCEDIR=.
SOURCES := $(find $(SOURCEDIR) -name '*.go')

# Binary info
BINARY := hermes
PKG_NAME := hermes
PKG_PREFIX := "github.com/TheHipbot/${PKG_NAME}"

GOOS := darwin
GOARCH := amd64

# Build info
VERSION := 1.0.0
BUILD := `git rev-parse --short HEAD`

.DEFAULT_GOAL: $(BINARY)

.PHONY: clean
clean:
	$(GOCMD) clean && if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

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
	GOOS=${GOOS} GOARCH=${GOARCH} $(GOCMD) build -o ${BINARY}

.PHONY: test
test: 
	go test ./...

.PHONY: docker-build
docker-build:
	docker build -t "thehipbot/${PKG_NAME}:${BUILD}" --build-arg PACKAGE=$(PKG_PREFIX) .