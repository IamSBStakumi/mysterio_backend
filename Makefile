init:
	go mod tidy
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
	oapi-codegen --config=docs/config.yaml docs/openapi.yaml