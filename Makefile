BINARY := nplus
PROJECT_ROOT := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
SOURCEDIR := .
SOURCE_FILES := $(shell find $(SOURCEDIR) -name '*.go')

.PHONY: all
all: pkg/darwin_amd64/$(BINARY)-darwin-amd64 pkg/linux_amd64/$(BINARY)-linux-amd64

pkg/darwin_amd64/$(BINARY)-darwin-amd64: $(SOURCE_FILES)
	GOOS=darwin GOARCH=amd64 \
	go build -v -o "$@"

pkg/linux_amd64/$(BINARY)-linux-amd64: $(SOURCE_FILES)
	GOOS=linux GOARCH=amd64 \
	go build -v -o "$@"

install:
	cp pkg/linux_amd64/$(BINARY)-linux-amd64 /usr/local/bin/$(BINARY)
	chmod 755 /usr/local/bin/$(BINARY)

uninstall:
	rm -f /usr/local/bin/$(BINARY)

.PHONY: deps
deps:
	go get -v -d

.PHONY: clean
clean:
	go clean -i -x -v
	rm -f pkg/darwin_amd64/$(BINARY)-darwin-amd64 pkg/linux_amd64/$(BINARY)-linux-amd64
