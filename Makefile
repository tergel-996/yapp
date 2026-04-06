VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X github.com/tergel/yapp/internal/cli.Version=$(VERSION)"

.PHONY: build test vet clean install install-app

build:
	go build $(LDFLAGS) -o bin/yapp-cli ./cmd/yapp

vet:
	go vet ./...

test:
	go test ./... -v

clean:
	rm -rf bin/

install: build
	cp bin/yapp-cli /usr/local/bin/yapp-cli

# Build yapp-cli and regenerate the Yapp.app bundle in ~/Applications.
# Handy during development for end-to-end testing of the .app launch path
# (Cmd+Tab visibility, terminal integration, etc.).
install-app: build
	./bin/yapp-cli install
