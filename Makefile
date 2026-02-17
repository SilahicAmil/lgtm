PREFIX ?= /usr/local

build:
	go build -o lgtm ./cmd

install: build
	install -d $(PREFIX)/bin
	install -m 755 lgtm $(PREFIX)/bin/lgtm
	install -d $(PREFIX)/share/lgtm
	install -m 644 config/patterns.txt $(PREFIX)/share/lgtm/patterns.txt

uninstall:
	rm -f $(PREFIX)/bin/lgtm
	rm -rf $(PREFIX)/share/lgtm

run:
	go run ./cmd

clean:
	rm -f lgtm