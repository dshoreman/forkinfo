PREFIX ?= /usr

all:
	@echo -n "Building forkinfo..."
	@go build -i -v && echo " [DONE]"

install:
	@echo "Preparing package structure"
	@mkdir -p "$(DESTDIR)$(PREFIX)/bin"

	@echo "Installing Forkinfo..."
	@mv forkinfo "$(DESTDIR)$(PREFIX)/bin/forkinfo"

uninstall:
	@echo "Uninstalling Forkinfo..."
	@rm -v "$(DESTDIR)$(PREFIX)/bin/forkinfo"
	@echo "Uninstall complete"
