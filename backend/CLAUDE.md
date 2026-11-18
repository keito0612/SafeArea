# SafeArea Backend - Claude Code 指示書

## プロジェクト概要

| 項目 | 内容 |
|------|------|
| アプリ名 | SafeArea Backend |
| 概要 | 災害リスク可視化APIサービス |
| 言語 | Go 1.22 |
| ランタイム | AWS Lambda |

## 技術スタック

| 項目 | 技術 |
|------|------|
| ルーター | chi |
| ORM | GORM |
| DB | PostgreSQL |
| バリデーション | go-playground/validator |
| ログ | log/slog (標準) |

## アーキテクチャ

レイヤードアーキテクチャを採用

```
Controller → UseCase → Repository → DB
                ↓
            外部API
```

### フォルダ構成

```
internal/
├── controller/           # コントローラー層（HTTP処理）
├── usecase/              # ユースケース層（ビジネスロジック）
│   ├── xxx_usecase.go    # インターフェース
│   └── xxx_usecase_impl.go # 実装
├── domain/               # ドメイン層
│   ├── model/            # エンティティ
│   └── repository/       # リポジトリインターフェース
├── infrastructure/       # インフラストラクチャ層
│   ├── persistence/      # リポジトリ実装
│   └── external/         # 外部API呼び出し
└── dto/                  # DTO
    ├── request/
    └── response/
```

## コーディング規約

### 命名規則

| 対象 | 規則 | 例 |
|------|------|-----|
| ファイル名 | snake_case | `risk_usecase.go` |
| 構造体 | PascalCase | `RiskUseCase` |
| インターフェース | PascalCase | `RiskRepository` |
| メソッド | PascalCase | `GetRisk` |
| 変数 | camelCase | `riskLevel` |

### インターフェース

- UseCaseとRepositoryは必ずインターフェースを定義
- テスト時にモック化できるようにする

```go
// Good
type RiskUseCase interface {
    GetRisk(ctx context.Context, address string) (*response.RiskResponse, error)
}

type riskUseCaseImpl struct {
    repo repository.RiskRepository
}

// Bad - インターフェースなし
type RiskUseCase struct {
    repo *RiskRepositoryImpl
}
```

### エラーハンドリング

```go
// Good - エラーをラップ
if err != nil {
    return nil, fmt.Errorf("failed to get risk: %w", err)
}

// Bad - エラー無視
risk, _ := repo.Find(ctx, id)
```

### コンテキスト

- 全てのメソッドの第一引数は `context.Context`

```go
// Good
func (u *riskUseCaseImpl) GetRisk(ctx context.Context, address string) (*RiskResponse, error)

// Bad
func (u *riskUseCaseImpl) GetRisk(address string) (*RiskResponse, error)
```

## エラーレスポンス形式

### 通常エラー

```json
{
  "success": false,
  "error": {
    "code": "ADDRESS_NOT_FOUND",
    "message": "住所が見つかりませんでした"
  }
}
```

### バリデーションエラー

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "入力内容に問題があります",
    "details": [
      { "field": "address", "message": "必須項目です" }
    ]
  }
}
```

### エラーコード一覧

| コード | HTTP | 意味 |
|--------|------|------|
| `VALIDATION_ERROR` | 400 | 入力値エラー |
| `ADDRESS_NOT_FOUND` | 404 | 住所から座標取得失敗 |
| `RISK_DATA_NOT_FOUND` | 404 | リスクデータなし |
| `EXTERNAL_API_ERROR` | 502 | 外部API失敗 |
| `DATABASE_ERROR` | 500 | DB失敗 |
| `INTERNAL_ERROR` | 500 | その他 |

## 外部API

| API | 用途 | タイムアウト | リトライ |
|-----|------|-------------|----------|
| 国土地理院 | ジオコーディング | 5秒 | 1回 |
| ハザードマップポータル | 災害データ | 10秒 | 1回 |
| J-SHIS | 地震・液状化 | 5秒 | 1回 |

## テスト

### テスト方針

- AAAパターン（Arrange-Act-Assert）
- モックはインターフェースで実装
- テストファイルは同じフォルダに配置

```go
func TestRiskUseCase_GetRisk(t *testing.T) {
    // Arrange
    mockRepo := &MockRiskRepository{}
    usecase := NewRiskUseCase(mockRepo)

    // Act
    result, err := usecase.GetRisk(ctx, "東京都渋谷区")

    // Assert
    assert.NoError(t, err)
    assert.Equal(t, "medium", result.OverallRisk)
}
```

### テスト実行

```bash
go test ./...
go test -cover ./...
go test -v ./...
```

## 開発コマンド

```bash
# ローカル実行
go run cmd/api/main.go

# DB起動
docker-compose up -d

# Lint
golangci-lint run

# ビルド（Lambda用）
GOOS=linux GOARCH=amd64 go build -o bootstrap cmd/api/main.go

# 災害データのインポート
go run scripts/import-disasters/main.go

# シード実行（災害データをDBに投入）
make seed-disasters
```

## 災害履歴データ

### データソース
- 内閣府・気象庁の公開データ（CSV/Excel）

### 取得・登録フロー
1. CSVダウンロード（内閣府・気象庁サイトから）
2. データ加工（`scripts/import-disasters/`）
3. DBに保存（シードスクリプト）

### 対象期間
- 過去50年（1975年〜現在）

詳細は [README.md](./README.md) を参照

## DB接続

- グローバル変数でコネクション再利用
- Lambda同時実行数: 10〜20に制限
- タイムアウト: 接続5秒、クエリ10秒

## 環境変数

| 変数名 | 説明 |
|--------|------|
| `ENV` | 環境 (local/dev/prod) |
| `DB_HOST` | DBホスト |
| `DB_PORT` | DBポート |
| `DB_NAME` | DB名 |
| `DB_USER` | DBユーザー |
| `DB_PASSWORD` | DBパスワード |
| `GEOCODING_API_URL` | 国土地理院API URL |
| `JSHIS_API_URL` | J-SHIS API URL |

## 参考ドキュメント

- [ARCHITECTURE.md](./ARCHITECTURE.md) - アーキテクチャ詳細
- [README.md](./README.md) - セットアップ手順
