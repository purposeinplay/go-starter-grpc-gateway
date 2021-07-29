CHECK_FILES?=$$(go list ./... | grep -v /vendor/)

.PHONY: proto
proto: ## Regenerate proto files.
	protoc -I=proto -I./vendor/github.com/grpc-ecosystem/grpc-gateway/v2 -I=apigrpc apigrpc/starter.proto --go_out=apigrpc --go-grpc_out=apigrpc --grpc-gateway_out=apigrpc --openapiv2_out=apigrpc

.PHONY: image
image: ## Build the Docker image.
	docker build -t go-starter .

.PHONY: test
test: ## Run tests.
	go test -p 1 -v $(CHECK_FILES)