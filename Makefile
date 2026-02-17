PREFIX ?= /usr/local

build:
	go build -o lgtm ./cmd

install: build
	install -d $(PREFIX)/bin
	install -m 755 lgtm $(PREFIX)/bin/lgtm

uninstall:
	rm -f $(PREFIX)/bin/lgtm

run:
	go run ./cmd

clean:
	rm -f lgtm