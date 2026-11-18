# SafeArea Backend Architecture

## 概要

| 項目 | 内容 |
|------|------|
| 言語 | Go 1.22 |
| ランタイム | AWS Lambda |
| ルーター | chi |
| ORM | GORM |
| DB | PostgreSQL (RDS) |

## フォルダ構成

```
backend/
├── cmd/
│   └── api/
│       └── main.go                     # エントリーポイント
│
├── internal/
│   ├── controller/                     # コントローラー層
│   │   ├── risk_controller.go
│   │   └── risk_controller_test.go
│   │
│   ├── usecase/                        # ユースケース層
│   │   ├── risk_usecase.go             # インターフェース定義
│   │   ├── risk_usecase_impl.go        # 実装
│   │   └── risk_usecase_test.go
│   │
│   ├── domain/                         # ドメイン層
│   │   ├── model/                      # エンティティ
│   │   │   └── risk.go
│   │   └── repository/                 # リポジトリインターフェース
│   │       └── risk_repository.go
│   │
│   ├── infrastructure/                 # インフラストラクチャ層
│   │   ├── persistence/                # リポジトリ実装
│   │   │   ├── risk_repository_impl.go
│   │   │   └── risk_repository_test.go
│   │   └── external/                   # 外部API
│   │       ├── geocoding.go
│   │       └── jshis.go
│   │
│   └── dto/                            # DTO
│       ├── request/
│       │   └── risk_request.go
│       └── response/
│           └── risk_response.go
│
├── docker-compose.yml                  # ローカル開発用DB
├── go.mod
├── go.sum
└── ARCHITECTURE.md                     # このファイル
```

## レイヤードアーキテクチャ

```
ユーザー
   │
   ▼ やりとり
┌─────────────────────────────────────────────────────┐
│              コントローラー層                        │
│         (Controller + DTO)                         │
│         HTTPリクエスト受付 / レスポンス返却            │
└─────────────────────────┬───────────────────────────┘
                          │ 依存（UseCaseインターフェース）
                          ▼
┌─────────────────────────────────────────────────────┐
│              ユースケース層                          │
│         (UseCase Interface + Impl)                 │
│              アプリケーションロジック                 │
└─────────────────────────┬───────────────────────────┘
                          │ 依存（Repositoryインターフェース）
                          ▼
┌─────────────────────────────────────────────────────┐
│              ドメイン層                              │
│         (Model + Repository Interface)             │
│         エンティティ / リポジトリインターフェース       │
└─────────────────────────────────────────────────────┘
                          ▲
                          │ 実装
┌─────────────────────────┴───────────────────────────┐
│              インフラストラクチャ層                   │
│    (Repository Impl + External API)                │
│         DB接続 / 外部API呼び出し                     │
└──────────────────┬──────────────────────────────────┘
                   │ やりとり
                   ▼
            ┌─────────────┐
            │ DB / 外部API │
            └─────────────┘
```

### インターフェースと実装の関係

```
┌─────────────┐      uses       ┌──────────────────────┐
│ Controller  │ ───────────────▶│ <<interface>>        │
└─────────────┘                 │ UseCase              │
                                └──────────────────────┘
                                           ▲
                                           │ implements
                                ┌──────────────────────┐
                                │ UseCaseImpl          │
                                └──────────┬───────────┘
                                           │ uses
                                           ▼
                                ┌──────────────────────┐
                                │ <<interface>>        │
                                │ Repository           │
                                └──────────────────────┘
                                           ▲
                                           │ implements
                                ┌──────────────────────┐
                                │ RepositoryImpl       │
                                └──────────────────────┘
```

### 各レイヤーの責務

| レイヤー | 責務 | 含むもの |
|----------|------|----------|
| コントローラー層 | HTTPリクエスト受付、バリデーション、レスポンス返却 | Controller, DTO |
| ユースケース層 | アプリケーションロジック | UseCase Interface + Impl |
| ドメイン層 | ビジネスエンティティ、リポジトリ契約 | Model, Repository Interface |
| インフラストラクチャ層 | 外部との通信 | Repository Impl, External API |

### インターフェース定義

#### UseCase インターフェース (usecase/risk_usecase.go)

```go
package usecase

import (
    "context"
    "github.com/isobekeito/safearea/internal/dto/response"
)

type RiskUseCase interface {
    GetRisk(ctx context.Context, address string) (*response.RiskResponse, error)
}
```

#### Repository インターフェース (domain/repository/risk_repository.go)

```go
package repository

import (
    "context"
    "github.com/isobekeito/safearea/internal/domain/model"
)

type RiskRepository interface {
    FindByLocation(ctx context.Context, lat, lng float64) (*model.Risk, error)
}

type DisasterRepository interface {
    FindByPrefecture(ctx context.Context, prefectureCode string) ([]*model.Disaster, error)
}
```

### テスト時のモック

```go
// Controller のテスト → UseCase をモック
type MockRiskUseCase struct {}
func (m *MockRiskUseCase) GetRisk(ctx context.Context, address string) (*response.RiskResponse, error) {
    return &response.RiskResponse{...}, nil
}

// UseCase のテスト → Repository をモック
type MockRiskRepository struct {}
func (m *MockRiskRepository) FindByLocation(ctx context.Context, lat, lng float64) (*model.Risk, error) {
    return &model.Risk{...}, nil
}
```

### AAAパターンによるテスト例

```go
func TestRiskUseCase_GetRisk(t *testing.T) {
    // Arrange（準備）
    mockRepo := &MockRiskRepository{}
    usecase := NewRiskUseCase(mockRepo)

    // Act（実行）
    result, err := usecase.GetRisk(context.Background(), "東京都渋谷区")

    // Assert（検証）
    assert.NoError(t, err)
    assert.Equal(t, "medium", result.OverallRisk)
}
```

### レイヤードの原則

1. **上位→下位のみ依存**: Controller → UseCase → Domain ← Infrastructure
2. **インターフェースに依存**: 具象クラスではなくインターフェースに依存
3. **レイヤーをスキップしない**: ControllerからRepositoryを直接呼ばない
4. **テスト容易性**: 各層をモック化してテスト可能

## エラーハンドリング

### エラーコード形式

- 英語 SNAKE_CASE

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

| コード | 意味 | HTTP |
|--------|------|------|
| `VALIDATION_ERROR` | 入力値エラー | 400 |
| `ADDRESS_NOT_FOUND` | 住所から座標取得失敗 | 404 |
| `RISK_DATA_NOT_FOUND` | リスクデータなし | 404 |
| `EXTERNAL_API_ERROR` | 外部API失敗 | 502 |
| `DATABASE_ERROR` | DB接続/クエリ失敗 | 500 |
| `INTERNAL_ERROR` | その他内部エラー | 500 |

## 外部API連携

### 使用API

| API | 用途 | タイムアウト |
|-----|------|-------------|
| 国土地理院 | 住所→緯度経度変換 | 5秒 |
| ハザードマップポータル | 洪水・土砂・津波 | 10秒 |
| J-SHIS | 地震・液状化 | 5秒 |

### リトライ設定

| 項目 | 値 |
|------|-----|
| リトライ回数 | 1回 |
| リトライ間隔 | 500ms |

### エラー時の挙動

- 外部APIエラー時は500エラーを返す
- エラーコード: `EXTERNAL_API_ERROR`

## DB接続

### コネクション管理

| 項目 | 設定 |
|------|------|
| 方式 | グローバル変数で再利用 |
| Lambda同時実行数 | 10〜20に制限 |

### タイムアウト設定

| 項目 | 値 |
|------|-----|
| 接続タイムアウト | 5秒 |
| クエリタイムアウト | 10秒 |
| アイドルタイムアウト | 5分 |

### 接続例

```go
var db *gorm.DB

func init() {
    var err error
    db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }

    sqlDB, _ := db.DB()
    sqlDB.SetMaxOpenConns(10)
    sqlDB.SetMaxIdleConns(5)
    sqlDB.SetConnMaxIdleTime(5 * time.Minute)
}
```

## 設定管理

### 方式

- 環境変数のみ（Private Subnetで保護されているため）

### 環境変数一覧

| 変数名 | 説明 | 例 |
|--------|------|-----|
| `ENV` | 環境 | `local` / `dev` / `prod` |
| `DB_HOST` | DBホスト | `localhost` / `xxx.rds.amazonaws.com` |
| `DB_PORT` | DBポート | `5432` |
| `DB_NAME` | DB名 | `safearea` |
| `DB_USER` | DBユーザー | `postgres` |
| `DB_PASSWORD` | DBパスワード | `****` |
| `GEOCODING_API_URL` | 国土地理院API | `https://msearch.gsi.go.jp/address-search/AddressSearch` |
| `JSHIS_API_URL` | J-SHIS API | `https://www.j-shis.bosai.go.jp/` |

### 環境分離

```
Lambda関数名:
  - safearea-api-dev   (ENV=dev)
  - safearea-api-prod  (ENV=prod)
```

### ローカル開発

```
.env.local        → 環境変数
docker-compose    → PostgreSQL起動
```

## テスト

### テスト方針

| 項目 | 方針 |
|------|------|
| 配置 | 同じフォルダに `_test.go` |
| モック | interfaceで実装 |
| 単体テスト | モックを使用 |
| 統合テスト | Docker DBを使用 |
| テスト手法 | AAA（Arrange-Act-Assert） |

### AAAパターン

```go
func TestRiskService_GetRisk(t *testing.T) {
    // Arrange（準備）
    mockRepo := &MockRiskRepository{}
    service := NewRiskService(mockRepo)

    // Act（実行）
    result, err := service.GetRisk(ctx, "東京都渋谷区")

    // Assert（検証）
    assert.NoError(t, err)
    assert.Equal(t, "medium", result.OverallRisk)
}
```

### テスト実行

```bash
# 単体テスト
go test ./...

# カバレッジ付き
go test -cover ./...

# 詳細表示
go test -v ./...
```

## 依存パッケージ

```go
require (
    // Lambda
    github.com/aws/aws-lambda-go
    github.com/awslabs/aws-lambda-go-api-proxy

    // Router
    github.com/go-chi/chi/v5
    github.com/go-chi/cors

    // ORM
    gorm.io/gorm
    gorm.io/driver/postgres

    // Validation
    github.com/go-playground/validator/v10

    // UUID
    github.com/google/uuid

    // Environment (local)
    github.com/joho/godotenv
)
```

## API エンドポイント

| Method | Path | 説明 |
|--------|------|------|
| GET | `/api/health` | ヘルスチェック |
| GET | `/api/risk?address=xxx` | 災害リスク取得 |
| GET | `/api/disasters?prefecture=xxx` | 過去の災害一覧 |

## レスポンス形式

### 成功時

```json
{
  "success": true,
  "data": {
    "overall_risk": "high",
    "risks": [
      {
        "type": "flood",
        "level": "high",
        "description": "浸水想定区域内"
      }
    ],
    "tags": ["flood", "liquefaction"],
    "disasters": []
  }
}
```

### 失敗時

```json
{
  "success": false,
  "error": {
    "code": "ADDRESS_NOT_FOUND",
    "message": "住所が見つかりませんでした"
  }
}
```
