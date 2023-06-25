ifeq ($(PREFIX),)
	PREFIX := /usr/local
endif

ifeq ($(BINDIR),)
	BINDIR := /bin
endif

CMD=ws_internal

make: ws_internal

ws_internal: main.go wrapper.go
	go build -o $@ $^

install:
	mkdir -p $(DESTDIR)$(PREFIX)$(BINDIR)/
	install -m 755 $(CMD) $(DESTDIR)$(PREFIX)$(BINDIR)/

uninstall:
	rm -rf $(DESTDIR)$(PREFIX)$(BINDIR)/$(CMD)

.PHONY: install uninstall make
