init:
	go mod tidy
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
	oapi-codegen --config=docs/config.yaml docs/openapi.yaml

# OpenAPIからコード生成
generate:
	@echo "Generating code from OpenAPI spec..."
	oapi-codegen -config docs/config.yaml docs/openapi.yaml
	@echo "Code generation complete!"
