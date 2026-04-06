VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X github.com/tergel/yapp/internal/cli.Version=$(VERSION)"

.PHONY: build test clean install

build:
	go build $(LDFLAGS) -o bin/yapp-cli ./cmd/yapp

test:
	go test ./... -v

clean:
	rm -rf bin/

install: build
	cp bin/yapp-cli /usr/local/bin/yapp-cli
