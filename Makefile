CHECK_FILES?=$$(go list ./... | grep -v /vendor/)

.PHONY: proto
proto: ## Regenerate proto files.
	protoc -I=proto -I./vendor/github.com/grpc-ecosystem/grpc-gateway/v2 -I=apigrpc apigrpc/v1/starter.proto --go_out=apigrpc/v1 --go-grpc_out=apigrpc/v1 --grpc-gateway_out=apigrpc/v1 --openapiv2_out=apigrpc

.PHONY: image
image: ## Build the Docker image.
	docker build -t go-starter .

test-ci:
	go test -mod=vendor -count=1 -timeout 300s -short -coverprofile=coverage.txt -covermode=atomic ./internal/...

.PHONY: test
test: ## Run tests.
	go test -p 1 -v $(CHECK_FILES)

.PHONY: migrate-test
migrate-test: ## Run migrations.
	go run main.go migrate --config config.test.yaml

.PHONY: check
lint:
	golangci-lint run --build-tags=dev