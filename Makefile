# ponysay-go Makefile

GO ?= go
PREFIX ?= /usr/local
BINDIR ?= $(PREFIX)/bin
BINARY_NAME ?= ponysay
SYMNAME ?= ponythink

.PHONY: all build test fmt vet check install uninstall clean

all: build

build:
	$(GO) build -o $(BINARY_NAME) ./cmd/ponysay
	ln -sf $(BINARY_NAME) $(SYMNAME)

test:
	$(GO) test -v ./...

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

check: fmt vet test

install: build
	mkdir -p $(DESTDIR)$(BINDIR)
	install -m 755 $(BINARY_NAME) $(DESTDIR)$(BINDIR)/$(BINARY_NAME)
	ln -sf $(BINARY_NAME) $(DESTDIR)$(BINDIR)/$(SYMNAME)

uninstall:
	rm -f $(DESTDIR)$(BINDIR)/$(BINARY_NAME)
	rm -f $(DESTDIR)$(BINDIR)/$(SYMNAME)

clean:
	rm -f $(BINARY_NAME) $(SYMNAME) main ponysay_test_bin ponythink_test_bin
