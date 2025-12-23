init:
	go mod tidy
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
	oapi-codegen --config=docs/config.yaml docs/openapi.yaml

# oapi-codegenを使えるようにPATHを通す
path:
	echo 'export PATH=$PATH:/home/sbs_takumi/go/bin' >> ~/.zshrc
	source ~/.zshrc

# ツールのインストール
install-tools:
	@echo "Installing oapi-codegen..."
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest

# OpenAPIからコード生成
generate:
	@echo "Generating code from OpenAPI spec..."
	oapi-codegen --config=docs/config.yaml docs/openapi.yaml
	@echo "Code generation complete!"

# サーバーを起動
run:
	go run cmd/server/main.go

# バイナリをビルド
build:
	go build -o server cmd/server/main.go

# テストを実行
test:
	go test -v ./...

# テストカバレッジを確認
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# 依存関係をダウンロード
deps:
	go mod download

# 依存関係を整理
tidy:
	go mod tidy

# OpenAPI仕様のバリデーション（オプション：spectral等が必要）
validate-openapi:
	@echo "Validating OpenAPI spec..."
	@command -v spectral >/dev/null 2>&1 || { echo "spectral not installed. Run: npm install -g @stoplight/spectral-cli"; exit 1; }
	spectral lint api/openapi.yaml

setup: install-tools deps generate
	@echo "Setup complete! Run 'make run' to start the server."

# ヘルプを表示
help:
	@echo "利用可能なコマンド:"
	@echo "  make install-tools  - 必要なツールをインストール"
	@echo "  make generate       - OpenAPIからコードを生成"
	@echo "  make run            - サーバーを起動"
	@echo "  make build          - バイナリをビルド"
	@echo "  make test           - テストを実行"
	@echo "  make coverage       - テストカバレッジを確認"
	@echo "  make deps           - 依存関係をダウンロード"
	@echo "  make tidy           - 依存関係を整理"
	@echo "  make validate-openapi - OpenAPI仕様をバリデーション"
	@echo "  make setup          - 初回セットアップ（ツール+依存関係+生成）"
	@echo "  make clean          - クリーンアップ"
