GOOS ?= $(shell go env GOOS)

.PHONY: server
server:
	GOOS=$(GOOS) CGO_ENABLED=0 go build -o _build/server -v ./cmd/server

.PHONY: client
client:
	GOOS=$(GOOS) CGO_ENABLED=0 go build -o _build/client -v ./cmd/client

.PHONY: build
build: server client

.PHONY: test
test:
	go test -v	./...

.PHONY: clean
clean:
	rm -fr _build/
