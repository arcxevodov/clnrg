.DEFAULT_GOAL := build

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
	cp clnrg /usr/bin
.PHONY:install

uninstall:
	rm /usr/bin/clnrg
.PHONY:uninstall