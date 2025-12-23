# Mysterio Backend

AI が生成するマーダーミステリーゲームのバックエンド API

## 技術スタック

- Go 1.21
- Echo Framework
- OpenAPI 3.0
- oapi-codegen（コード自動生成）
- インメモリストレージ（MVP 版）

## ディレクトリ構成

```
mysterio_backend/
├── api/
│   ├── openapi.yaml         # OpenAPI仕様書
│   ├── oapi-codegen.yaml    # コード生成設定
│   └── generated.go         # 自動生成されたコード
├── cmd/
│   └── server/              # メインエントリーポイント
├── internal/
│   ├── handler/             # HTTPハンドラ
│   ├── service/             # ビジネスロジック
│   ├── repository/          # データアクセス層
│   ├── domain/              # ドメインモデル
│   └── ai/                  # AIシナリオ生成
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── go.mod
```

## セットアップ

### 前提条件

- Go 1.21 以上
- Docker & Docker Compose（オプション）

### 初回セットアップ

```bash
# 1. リポジトリをクローン
git clone <repository-url>
cd mysterio_backend

# 2. 必要なツールのインストールと依存関係の解決、コード生成を一括実行
make setup

# または個別に実行
make install-tools  # oapi-codegenのインストール
make deps           # 依存関係のダウンロード
make generate       # コード生成
```

### コード生成

OpenAPI 仕様書（`api/openapi.yaml`）を更新したら、以下のコマンドでコードを再生成します：

```bash
make generate
```

このコマンドは以下を実行します：

- `api/openapi.yaml` から Go のコードを自動生成
- `api/generated.go` に出力（型定義、インターフェース、ルーティング）

### サーバーの起動

**方法 A: Go 直接実行**

```bash
make run
```

**方法 B: Makefile を使用**

```bash
make run
```

**方法 C: Docker を使用**

```bash
make docker-up
```

サーバーは `http://localhost:8080` で起動します。

## API ドキュメント

### エンドポイント一覧

すべてのエンドポイントは `/api/v1` プレフィックスが付きます。

#### 1. セッション作成

```http
POST /api/v1/sessions
Content-Type: application/json

{
  "playerCount": 3,
  "difficulty": "easy"
}
```

**レスポンス例**

```json
{
  "sessionId": "550e8400-e29b-41d4-a716-446655440000",
  "ownerPlayerId": "660e8400-e29b-41d4-a716-446655440001",
  "playerIds": [
    "660e8400-e29b-41d4-a716-446655440001",
    "660e8400-e29b-41d4-a716-446655440002",
    "660e8400-e29b-41d4-a716-446655440003"
  ],
  "initialPhase": {
    "phaseNumber": 0,
    "phaseType": "introduction",
    "description": "事件の概要と各自の立場を確認",
    "publicText": "被害者：サミュエル・ゴールドスタイン...",
    "duration": 10
  }
}
```

#### 2. フェーズ情報取得

```http
GET /api/v1/sessions/{sessionId}/phase
X-Player-Id: {playerId}
```

**レスポンス例**

```json
{
  "phaseNumber": 0,
  "phaseType": "introduction",
  "description": "事件の概要と各自の立場を確認",
  "publicText": "被害者：サミュエル・ゴールドスタイン...",
  "privateText": "【あなたの役割】\n名前: アレックス・クロフォード...",
  "availableActions": ["advance_phase"]
}
```

#### 3. フェーズ進行

```http
POST /api/v1/sessions/{sessionId}/phase/advance
X-Player-Id: {ownerPlayerId}
```

**レスポンス例**

```json
{
  "message": "phase advanced successfully"
}
```

#### 4. 投票

```http
POST /api/v1/sessions/{sessionId}/vote
X-Player-Id: {playerId}
Content-Type: application/json

{
  "targetPlayerId": "660e8400-e29b-41d4-a716-446655440002"
}
```

**レスポンス例**

```json
{
  "message": "vote recorded successfully"
}
```

#### 5. 結果取得

```http
GET /api/v1/sessions/{sessionId}/result
```

**レスポンス例**

```json
{
  "truth": "犯人はアレックス・クロフォードである。動機は...",
  "winnerId": "660e8400-e29b-41d4-a716-446655440001",
  "winnerName": "アレックス・クロフォード",
  "votingResult": {
    "660e8400-e29b-41d4-a716-446655440001": 2,
    "660e8400-e29b-41d4-a716-446655440002": 1
  },
  "playerNames": {
    "660e8400-e29b-41d4-a716-446655440001": "アレックス・クロフォード",
    "660e8400-e29b-41d4-a716-446655440002": "エミリー・ハートウェル"
  }
}
```

### エラーレスポンス

すべてのエラーは以下の形式で返されます：

```json
{
  "message": "error description"
}
```

ステータスコード：

- `400 Bad Request`: リクエストが不正
- `403 Forbidden`: 権限がない
- `404 Not Found`: リソースが見つからない
- `500 Internal Server Error`: サーバーエラー

## 動作確認

### curl を使用したテスト

```bash
# 1. セッション作成
SESSION_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/sessions \
  -H "Content-Type: application/json" \
  -d '{"playerCount": 3, "difficulty": "easy"}')

echo $SESSION_RESPONSE | jq

# 2. sessionIdとplayerIdsを抽出
SESSION_ID=$(echo $SESSION_RESPONSE | jq -r '.sessionId')
OWNER_ID=$(echo $SESSION_RESPONSE | jq -r '.ownerPlayerId')
PLAYER_2=$(echo $SESSION_RESPONSE | jq -r '.playerIds[1]')

# 3. フェーズ情報を取得（プレイヤー1）
curl -H "X-Player-Id: $OWNER_ID" \
  http://localhost:8080/api/v1/sessions/$SESSION_ID/phase | jq

# 4. フェーズを進める
curl -X POST \
  -H "X-Player-Id: $OWNER_ID" \
  http://localhost:8080/api/v1/sessions/$SESSION_ID/phase/advance | jq

# 5. 投票フェーズまで進めて投票
# (フェーズを複数回進める)
curl -X POST \
  -H "X-Player-Id: $OWNER_ID" \
  http://localhost:8080/api/v1/sessions/$SESSION_ID/phase/advance

curl -X POST \
  -H "X-Player-Id: $OWNER_ID" \
  http://localhost:8080/api/v1/sessions/$SESSION_ID/phase/advance

curl -X POST \
  -H "X-Player-Id: $OWNER_ID" \
  http://localhost:8080/api/v1/sessions/$SESSION_ID/phase/advance

# 投票
curl -X POST \
  -H "X-Player-Id: $OWNER_ID" \
  -H "Content-Type: application/json" \
  -d "{\"targetPlayerId\": \"$PLAYER_2\"}" \
  http://localhost:8080/api/v1/sessions/$SESSION_ID/vote | jq

# 6. 結果を取得
curl http://localhost:8080/api/v1/sessions/$SESSION_ID/result | jq
```

## 開発ワークフロー

### OpenAPI 仕様の更新

1. `api/openapi.yaml` を編集
2. コードを再生成: `make generate`
3. 必要に応じて `internal/handler/api_handler.go` を更新
4. テストして確認: `make run`

### OpenAPI 仕様のバリデーション（オプション）

Spectral を使用して OpenAPI 仕様をバリデーションできます：

```bash
# Spectralのインストール（Node.js必要）
npm install -g @stoplight/spectral-cli

# バリデーション実行
make validate-openapi
```

## AWS デプロイ

### ECS Fargate へのデプロイ

#### 1. ECR リポジトリの作成

```bash
aws ecr create-repository --repository-name mysterio-backend --region ap-northeast-1
```

#### 2. イメージのビルドとプッシュ

```bash
# ECRにログイン
aws ecr get-login-password --region ap-northeast-1 | docker login --username AWS --password-stdin {account-id}.dkr.ecr.ap-northeast-1.amazonaws.com

# イメージをビルド
make docker-build

# タグ付け
docker tag mysterio-backend:latest {account-id}.dkr.ecr.ap-northeast-1.amazonaws.com/mysterio-backend:latest

# プッシュ
docker push {account-id}.dkr.ecr.ap-northeast-1.amazonaws.com/mysterio-backend:latest
```

#### 3. ECS クラスタとサービスの作成

ECS コンソール、または Terraform/CDK を使用してクラスタとサービスを作成します。

**タスク定義の設定例:**

- CPU: 256 (.25 vCPU)
- メモリ: 512 MB
- ポートマッピング: 8080
- 環境変数: PORT=8080

## フロントエンド連携

### TypeScript 型定義の生成

フロントエンド（Next.js）用に TypeScript 型定義を生成できます：

```bash
# openapi-typescriptをインストール（フロントエンドプロジェクトで）
npm install -D openapi-typescript

# 型定義を生成
npx openapi-typescript http://localhost:8080/api/openapi.yaml -o types/api.ts
```

または、OpenAPI 仕様ファイルから直接：

```bash
npx openapi-typescript ../mysterio_backend/api/openapi.yaml -o types/api.ts
```

## Make コマンド一覧

```bash
make install-tools   # 必要なツールをインストール
make generate        # OpenAPIからコードを生成
make run             # サーバーを起動
make build           # バイナリをビルド
make test            # テストを実行
make coverage        # テストカバレッジを確認
make deps            # 依存関係をダウンロード
make tidy            # 依存関係を整理
make validate-openapi # OpenAPI仕様をバリデーション
make setup           # 初回セットアップ（全て実行）
make clean           # クリーンアップ
make help            # ヘルプを表示
```

## トラブルシューティング

### コード生成エラー

```bash
# oapi-codegenが見つからない場合
make install-tools

# PATHに追加（~/.bashrc or ~/.zshrc）
export PATH=$PATH:$(go env GOPATH)/bin
```

### ポートが既に使用されている

```bash
# 使用中のプロセスを確認
lsof -i :8080

# プロセスを終了
kill -9 <PID>
```

## 今後の実装予定

- [ ] Anthropic API 統合（実際の AI シナリオ生成）
- [ ] Redis/PostgreSQL へのデータストア移行
- [ ] セッションの TTL 管理
- [ ] ログ強化
- [ ] メトリクス収集
- [ ] 認証・認可の強化
- [ ] WebSocket 対応（リアルタイム通知）
