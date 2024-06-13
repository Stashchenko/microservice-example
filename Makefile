OUTPUT := bin/service

GOOSE_CMD := goose -dir internal/migrations postgres "user=postgres dbname=go_grpc_example sslmode=disable"

goose-status:
	$(GOOSE_CMD) status
goose-up:
	$(GOOSE_CMD) up
goose-down:
	$(GOOSE_CMD) down

.PHONY: lint
lint:
	golangci-lint --enable gosec,misspell run ./...

.PHONY: build
build: lint
	go build -o $(OUTPUT) github.com/stashchenko/microservice-example/cmd

.PHONY: test
test:
	go test -v --race ./...

.PHONY: govulncheck
govulncheck:
	govulncheck ./...

.PHONY: gogen
gogen:
	go generate ./...

.PHONY: proto-health
proto-health:
	protoc \
	--go_out=. \
	--go_opt=paths=source_relative \
    --go-grpc_out=. \
	--go-grpc_opt=paths=source_relative \
    pkg/proto/health.proto

.PHONY: proto-all
proto-all:
	protoc \
	--go_out=. \
	--go_opt=paths=source_relative \
    --go-grpc_out=. \
	--go-grpc_opt=paths=source_relative \
    pkg/proto/*.proto
