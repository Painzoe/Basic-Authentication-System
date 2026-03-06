BINARY_NAME=auth-server
PREFIX=/usr/local/bin

all: build

build:
	go build -ldflags="-s -w" -o $(BINARY_NAME) main.go

clean:
	rm -f $(BINARY_NAME)
	rm -f auth.db

install: build
	install -m 755 $(BINARY_NAME) $(PREFIX)/$(BINARY_NAME)

uninstall:
	rm -f $(PREFIX)/$(BINARY_NAME)

test:
	go test ./...

.PHONY: all build clean install uninstall test
