VERSION ?= $(shell git describe --tags --dirty)
NAME := yk
GOFILES = *.go
# Multi-arch targets are generated from this
TARGET_ALIAS = $(NAME)-linux-amd64 $(NAME)-linux-arm $(NAME)-linux-arm64 $(NAME)-darwin-amd64
TARGETS = $(addprefix dist/,$(TARGET_ALIAS))

.PHONY: default
default: yk

.PHONY: all
all: $(TARGETS)

.PHONY: test
test: check
	go test .

.PHONY: build
build: build/yk

build/yk: $(GOFILES)
	mkdir -p build
	go build -o build/yk

.PHONY: clean
clean:
	rm -fr build/
	rm -fr dist/

.PHONY: check
check:
	pre-commit run --all-files

.PHONY: install-hooks
install-hooks:
	pre-commit install -f --install-hooks

# Distribution targets
$(TARGETS): $(GOFILES)
	mkdir -p ./dist
	GOOS=$(word 2, $(subst -, ,$(@))) GOARCH=$(word 3, $(subst -, ,$(@))) \
		 go build -ldflags '-X "main.version=${VERSION}"' -a \
		 -o $@

.PHONY: $(TARGET_ALIAS)
$(TARGET_ALIAS):
	$(MAKE) $(addprefix dist/,$@)
