PREFIX ?= /usr

install:
	@go build -o jota main.go
	@install -Dm755 jota $(DESTDIR)$(PREFIX)/bin/jota

uninstall:
	@rm -rf $(DESTDIR)$(PREFIX)/bin/jota