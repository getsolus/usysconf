.POSIX:
.SUFFIXES:
.SUFFIXES: .1 .5 .1.scd .5.scd

PKGNAME=usysconf

VERSION=0.6.0

VPATH=doc
PREFIX?=/usr/local
BINDIR?=$(DESTDIR)$(PREFIX)/bin
SYSDIR?=$(DESTDIR)/etc/$(PKGNAME).d
USRDIR?=$(DESTDIR)$(PREFIX)/share/default/$(PKGNAME).d
LOGDIR?=$(DESTDIR)/var/log/$(PKGNAME)
MANDIR?=$(DESTDIR)$(PREFIX)/share/man
GO?=go
GOFLAGS?=

GOSRC!=find . -name '*.go'
GOSRC+=go.mod go.sum

usysconf: $(GOSRC)
	$(GO) build $(GOFLAGS) \
		-ldflags "-X main.Prefix=$(PREFIX) \
		-X main.Version=$(VERSION) \
		-X main.LogDir=$(LOGDIR) \
		-X main.SysDir=$(SYSDIR) \
		-X main.UsrDir=$(USRDIR)" \
		-o $@

DOCS := \
	$(PKGNAME).1 \
	$(PKGNAME)-run.1 \
	$(PKGNAME)-config.5

.1.scd.1:
	scdoc < $< > $@

.5.scd.5:
	scdoc < $< > $@

doc: $(DOCS)

all: usysconf doc 

# Exists in GNUMake but not in NetBSD make and others.
RM?=rm -f

clean:
	$(GO) mod tidy
	$(RM) $(DOCS) $(PKGNAME) *.tar.gz
	$(RM) -r vendor

install: all
	mkdir -m755 -p $(BINDIR) $(USRDIR) $(SYSDIR) $(LOGDIR) $(MANDIR)/man1 \
		$(MANDIR)/man5
	install -m755 $(PKGNAME) $(BINDIR)/$(PKGNAME)
	install -m644 $(PKGNAME).1 $(MANDIR)/man1/$(PKGNAME).1
	install -m644 $(PKGNAME)-run.1 $(MANDIR)/man1/$(PKGNAME)-run.1
	install -m644 $(PKGNAME)-config.5 $(MANDIR)/man5/$(PKGNAME)-config.5

RMDIR_IF_EMPTY:=sh -c '\
if test -d $$0 && ! ls -1qA $$0 | grep -q . ; then \
	rmdir $$0; \
fi'

uninstall:
	$(RM) $(BINDIR)/$(PKGNAME)
	$(RM) $(MANDIR)/man1/$(PKGNAME).1
	$(RM) $(MANDIR)/man5/$(PKGNAME).5
	$(RM) -r $(LOGDIR)
	$(RM) -r $(SYSDIR)
	$(RM) -r $(USRDIR)
	$(RMDIR_IF_EMPTY) $(BINDIR)
	$(RMDIR_IF_EMPTY) $(MANDIR)/man1
	$(RMDIR_IF_EMPTY) $(MANDIR)/man5
	$(RMDIR_IF_EMPTY) $(MANDIR)

check:
	$(GO) get -u golang.org/x/lint/golint
	$(GO) get -u github.com/securego/gosec/cmd/gosec
	$(GO) get -u honnef.co/go/tools/cmd/staticcheck
	$(GO) get -u gitlab.com/opennota/check/cmd/aligncheck
	$(GO) fmt -x ./...
	$(GO) vet ./...
	golint -set_exit_status `go list ./... | grep -v vendor`
	gosec -exclude=G204 ./...
	staticcheck ./...
	aligncheck ./...
	$(GO) test -cover ./...

vendor: check clean
	$(GO) mod vendor

package: vendor
	tar --exclude='.git' \
		--exclude='*.tar.gz' \
	       	--exclude='examples' \
	       	--exclude="tags" \
	       	--exclude=".vscode" \
		--exclude="*~" \
		-zcvf $(PKGNAME)-v$(VERSION).tar.gz ../$(PKGNAME)

.DEFAULT_GOAL := all

.PHONY: all doc clean install uninstall check vendor package
