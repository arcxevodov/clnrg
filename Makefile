.DEFAULT_GOAL := build
SHELL := /bin/bash

fmt:
	go fmt ./...
.PHONY:fmt

vet: fmt
	go install golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow@latest
	go vet ./...
	shadow ./...
.PHONY:vet

build: vet
	go build -o clnrg -ldflags "-w -s" main.go
.PHONY:build

install:
	@if [[ ! -f /usr/local/clnr/clnr ]]; then \
		echo -e "\033[0;31mExecutable file clnr not found. Please run install_clnr.sh "; \
	else \
		mkdir /usr/local/clnrg; \
		cp clnrg /usr/local/clnrg; \
		ln -sf /usr/local/clnrg/clnrg /usr/local/sbin/clnrg; \
		ln -sf /usr/local/clnrg/clnrg /usr/local/bin/clnrg; \
	fi
.PHONY:install

uninstall:
	rm -r /usr/local/clnrg
	rm /usr/local/sbin/clnrg
	rm /usr/local/bin/clnrg
.PHONY:uninstall