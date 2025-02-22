BINARY := ./bin/bin

.PHONY: build
build:
	go build -o $(BINARY) ./cmd/conductor/main.go

.PHONY: run
run: build
	$(BINARY)
