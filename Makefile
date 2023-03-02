# Hey there! This is the makefile, meant specifically for Linux and macOS.


# Install locally to your system (Linux & macOS)
install:
	@go build -ldflags "-s -w" -o jota
	@install -Dm755 jota $(DESTDIR)/usr/local/bin/jota

uninstall:
	@rm -rf $(DESTDIR)/usr/local/bin/jota

# Remove the binary
clean:
	@rm -f jota




# Don't touch! Meant for package managers!
global_install:
	@go build -ldflags "-s -w" -o jota
	@install -Dm755 jota $(DESTDIR)/usr/bin/jota

global_uninstall:
	@rm -rf $(DESTDIR)/usr/bin/jota
