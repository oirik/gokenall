NAME      := $(shell basename `pwd`)
VERSION   := $(shell git describe --tags --abbrev=0)
REVISION  := $(shell git rev-parse --short HEAD)
GODEP     := $(shell command -v dep 2> /dev/null)
GOLINT    := $(shell command -v golint 2> /dev/null)
LDFLAGS   := -X 'main.Version=$(VERSION)' -X 'main.Revision=$(REVISION)'
DISTDIR   :=./dist
VENDORDIR :=./vendor
DIST_DIRS := find * -type d -exec

.PHONY: test
test: lint
	go test -race -v ./...

.PHONY: godep
godep:
ifndef GODEP
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
endif

.PHONY: golint
golint:
ifndef GOLINT
	go get -u github.com/golang/lint/golint
endif

.PHONY: deps
deps: godep
	dep ensure

.PHONY: build
build: deps
	go build -ldflags "$(LDFLAGS)" ./...

.PHONY: clean
clean:
	go clean
	rm -rf $(DISTDIR)/*
	rm -rf $(VENDORDIR)/*

.PHONY: lint
lint: golint deps
	go vet ./...
	golint -set_exit_status ./...

.PHONY: install
install: test
	go install -ldflags "$(LDFLAGS)"

.PHONY: cross-build
cross-build: test
	rm -rf $(DISTDIR)/*
	for os in darwin linux windows; do \
		for arch in amd64 386; do \
			GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 go build -a -ldflags "$(LDFLAGS)" -o dist/$$os-$$arch/$(NAME); \
			if [ "$${os}" = "windows" ]; then \
				mv dist/$$os-$$arch/$(NAME) dist/$$os-$$arch/$(NAME).exe; \
			fi; \
		done; \
	done

.PHONY: dist
dist: cross-build
	cd dist && \
	$(DIST_DIRS) cp ../LICENSE {} \; && \
	$(DIST_DIRS) cp ../README.md {} \; && \
	$(DIST_DIRS) tar -zcf $(NAME)-${VERSION}-{}.tar.gz {} \; && \
	$(DIST_DIRS) zip -r $(NAME)-${VERSION}-{}.zip {} \; && \
	cd ..
