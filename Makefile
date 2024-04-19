CGO_CPPFLAGS ?= ${CPPFLAGS}
export CGO_CPPFLAGS
CGO_CFLAGS ?= ${CFLAGS}
export CGO_CFLAGS
CGO_LDFLAGS ?= $(filter -g -L% -l% -O%,${LDFLAGS})
export CGO_LDFLAGS

EXE =
ifeq ($(shell go env GOOS),windows)
EXE = .exe
endif

## The following tasks delegate to `script/build.go` so they can be run cross-platform.

.PHONY: bin/pinecone$(EXE)
bin/pinecone$(EXE): script/build$(EXE)
	@script/build$(EXE) $@

script/build$(EXE): script/build.go
ifeq ($(EXE),)
	GOOS= GOARCH= GOARM= GOFLAGS= CGO_ENABLED= go build -o $@ $<
else
	go build -o $@ $<
endif

.PHONY: clean
clean: script/build$(EXE)
	@$< $@

.PHONY: completions
completions: bin/pinecone$(EXE)
	mkdir -p ./share/bash-completion/completions ./share/fish/vendor_completions.d ./share/zsh/site-functions
	bin/pinecone$(EXE) completion -s bash > ./share/bash-completion/completions/pinecone
	bin/pinecone$(EXE) completion -s fish > ./share/fish/vendor_completions.d/pinecone.fish
	bin/pinecone$(EXE) completion -s zsh > ./share/zsh/site-functions/_pinecone

# just a convenience task around `go test`
.PHONY: test
test:
	go test ./...

DESTDIR :=
prefix  := /usr/local
bindir  := ${prefix}/bin
datadir := ${prefix}/share
mandir  := ${datadir}/man

.PHONY: install
install: bin/pinecone # completions
	install -d ${DESTDIR}${bindir}
	install -m755 bin/pinecone ${DESTDIR}${bindir}/
	# install -d ${DESTDIR}${mandir}/man1
	# install -m644 ./share/man/man1/* ${DESTDIR}${mandir}/man1/
	# install -d ${DESTDIR}${datadir}/bash-completion/completions
	# install -m644 ./share/bash-completion/completions/pinecone ${DESTDIR}${datadir}/bash-completion/completions/pinecone
	# install -d ${DESTDIR}${datadir}/fish/vendor_completions.d
	# install -m644 ./share/fish/vendor_completions.d/pinecone.fish ${DESTDIR}${datadir}/fish/vendor_completions.d/pinecone.fish
	# install -d ${DESTDIR}${datadir}/zsh/site-functions
	# install -m644 ./share/zsh/site-functions/_pinecone ${DESTDIR}${datadir}/zsh/site-functions/_pinecone

.PHONY: uninstall
uninstall:
	rm -f ${DESTDIR}${bindir}/pinecone ${DESTDIR}${mandir}/man1/pinecone.1 ${DESTDIR}${mandir}/man1/pinecone-*.1
	rm -f ${DESTDIR}${datadir}/bash-completion/completions/pinecone
	rm -f ${DESTDIR}${datadir}/fish/vendor_completions.d/pinecone.fish
	rm -f ${DESTDIR}${datadir}/zsh/site-functions/_pinecone