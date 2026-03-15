ENTRY_POINT = cmd/api/main.go

.PHONY: server-build
server-build:
	@go build -o tmp/bel-server ${ENTRY_POINT}

.PHONY: server-run
server-run:
	@go run ${ENTRY_POINT}

.PHONY: clean
clean:
	@echo "cleaning tmp dir..."

	@rm -rf tmp

.PHONY: watch
watch:
	@air

.PHONY: server-dev
server-dev: clean watch
