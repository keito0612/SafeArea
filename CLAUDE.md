# SafeArea - プロジェクト指示書

## プロジェクト概要

| 項目 | 内容 |
|------|------|
| アプリ名 | SafeArea |
| コンセプト | 引っ越し前に災害リスクがわかる地図サービス |
| ターゲット | 引っ越し検討者、物件購入検討者 |
| 対象地域 | 全国（日本） |

## 技術スタック

### フロントエンド

- **フレームワーク**: Next.js 14 (App Router)
- **言語**: TypeScript
- **アーキテクチャ**: レイヤードアーキテクチャ
- **スタイリング**: Tailwind CSS
- **状態管理**: useState + Context
- **データフェッチ**: TanStack Query
- **フォーム**: React Hook Form
- **バリデーション**: Zod
- **地図**: Leaflet / React-Leaflet
- **テスト**: Vitest / React Testing Library / Playwright

### バックエンド

- **言語**: Go
- **ランタイム**: AWS Lambda
- **API**: API Gateway (HTTP API)
- **ルーティング**: aws-lambda-go + chi
- **ORM**: GORM
- **データベース**: PostgreSQL (RDS)

### インフラ

- **IaC**: Terraform
- **ホスティング**: AWS Amplify Hosting
- **DB**: RDS PostgreSQL (t4g.micro)
- **NAT**: NAT Instance (t4g.nano)
- **CDN**: CloudFront

## フォルダ構成

```
SafeArea/
├── frontend/                    # Next.js（レイヤードアーキテクチャ）
│   ├── app/                    # UI層（ページ）
│   │   ├── page.tsx            # LP
│   │   ├── search/
│   │   │   ├── page.tsx        # 検索ページ
│   │   │   └── results/
│   │   │       └── page.tsx    # 検索結果ページ
│   │   ├── terms/
│   │   │   └── page.tsx        # 利用規約
│   │   └── privacy/
│   │       └── page.tsx        # プライバシーポリシー
│   ├── components/             # UI層（コンポーネント）
│   │   ├── ui/                 # 共通UI
│   │   ├── map/                # 地図関連
│   │   ├── search/             # 検索関連
│   │   ├── risk/               # リスク表示関連
│   │   └── layout/             # レイアウト
│   ├── hooks/                  # Hooks層
│   ├── services/               # Services層
│   ├── repositories/           # Repositories層
│   ├── types/                  # 型定義
│   ├── lib/                    # ユーティリティ
│   ├── __tests__/              # テスト
│   └── public/
│
├── backend/                     # Go Lambda（クリーンアーキテクチャ）
│   ├── cmd/
│   │   └── api/
│   │       └── main.go         # Lambdaエントリーポイント
│   ├── internal/
│   │   ├── controller/         # コントローラー層（HTTP処理）
│   │   ├── usecase/            # ユースケース層（ビジネスロジック）
│   │   ├── domain/             # ドメイン層
│   │   │   ├── model/          # エンティティ
│   │   │   └── repository/     # リポジトリインターフェース
│   │   ├── infrastructure/     # インフラストラクチャ層
│   │   │   ├── persistence/    # リポジトリ実装
│   │   │   └── external/       # 外部API呼び出し
│   │   └── dto/                # DTO
│   ├── db/
│   │   ├── migrations/         # DBマイグレーション
│   │   └── seeds/              # 初期データ（災害履歴CSV等）
│   ├── scripts/
│   │   └── import-disasters/   # 災害データCSV加工スクリプト
│   ├── go.mod
│   └── go.sum
│
├── infrastructure/              # Terraform
│   ├── main.tf
│   ├── variables.tf
│   ├── outputs.tf
│   ├── vpc.tf
│   ├── security_groups.tf
│   ├── rds.tf
│   ├── lambda.tf
│   ├── api_gateway.tf
│   ├── nat_instance.tf
│   └── terraform.tfvars        # ※gitignore対象
│
├── docs/
│   ├── architecture.drawio     # インフラ構成図
│   └── api.md                  # API仕様
│
├── CLAUDE.md                   # このファイル
└── README.md
```

## コーディング規約

### 共通

- コミットメッセージは日本語OK
- ブランチ名: `feature/xxx`, `fix/xxx`, `docs/xxx`
- PRは日本語で記述

### TypeScript (フロントエンド)

- ESLint + Prettier を使用
- 関数コンポーネントを使用（クラスコンポーネント禁止）
- 型は明示的に定義（any禁止）
- ファイル名: kebab-case (`risk-card.tsx`)
- コンポーネント名: PascalCase (`RiskCard`)

```typescript
// Good
const RiskCard: React.FC<RiskCardProps> = ({ risk }) => {
  return <div>{risk.name}</div>;
};

// Bad
function RiskCard(props: any) {
  return <div>{props.risk.name}</div>;
}
```

### Go (バックエンド)

- gofmt + golangci-lint を使用
- エラーは必ずハンドリング
- 構造体タグは `json:"xxx"` 形式
- ファイル名: snake_case (`risk_service.go`)

```go
// Good
func (s *RiskService) GetRisk(ctx context.Context, address string) (*Risk, error) {
    risk, err := s.repo.FindByAddress(ctx, address)
    if err != nil {
        return nil, fmt.Errorf("failed to get risk: %w", err)
    }
    return risk, nil
}

// Bad
func (s *RiskService) GetRisk(address string) *Risk {
    risk, _ := s.repo.FindByAddress(address) // エラー無視
    return risk
}
```

### Terraform

- リソース名: snake_case
- 変数は `variables.tf` に集約
- 機密情報は `terraform.tfvars` に記載し、gitignore対象

## API設計

### エンドポイント

| Method | Path | 説明 |
|--------|------|------|
| GET | `/api/health` | ヘルスチェック |
| GET | `/api/search?address=xxx` | 災害リスク検索（詳細タブ用） |
| GET | `/api/search/disasters?address=xxx` | 過去の災害履歴（災害履歴タブ用） |

### レスポンス形式

#### GET `/api/search?address=xxx`

```json
{
  "success": true,
  "data": {
    "address": "東京都江東区豊洲1丁目",
    "coordinates": { "lat": 35.655, "lng": 139.793 },
    "overall_risk": "high",
    "risks": [
      { "type": "flood", "level": "high", "description": "浸水想定区域内" }
    ],
    "tags": ["flood", "liquefaction"],
    "disaster_history_count": 3
  }
}
```

#### GET `/api/search/disasters?address=xxx`

```json
{
  "success": true,
  "data": {
    "disasters": [
      {
        "id": "d001",
        "type": "earthquake",
        "name": "東日本大震災",
        "year": 2011,
        "damage": "液状化により道路陥没、断水被害"
      }
    ]
  }
}
```

※ 過去50年（1975年〜現在）のデータを返却
※ データソース：内閣府・気象庁の公開データ（初期データとしてDB登録）

## データベース設計

### テーブル一覧

- `disasters` - 過去の災害データ
- `disaster_areas` - 災害の被災地域
- `earthquake_risks` - 地震危険度データ（メッシュ単位）

### マイグレーション

```bash
# backend/db/migrations/ 配下に作成
# ファイル名: YYYYMMDDHHMMSS_description.sql
```

## 外部API

| API | 用途 | ドキュメント |
|-----|------|-------------|
| 国土地理院 | ジオコーディング | https://msearch.gsi.go.jp/address-search/AddressSearch |
| ハザードマップポータル | 洪水・土砂・津波タイル | https://disaportal.gsi.go.jp/ |
| J-SHIS | 地震・液状化データ | https://www.j-shis.bosai.go.jp/ |

## 開発コマンド

### フロントエンド

```bash
cd frontend

# 開発サーバー起動
npm run dev

# ビルド
npm run build

# Lint
npm run lint

# 型チェック
npm run type-check

# テスト
npm run test

# E2Eテスト
npm run test:e2e
```

### バックエンド

```bash
cd backend

# ローカル実行
go run cmd/api/main.go

# ビルド（Lambda用）
GOOS=linux GOARCH=amd64 go build -o bootstrap cmd/api/main.go

# テスト
go test ./...

# Lint
golangci-lint run
```

### Terraform

```bash
cd infrastructure

# 初期化
terraform init

# プラン確認
terraform plan

# 適用
terraform apply

# 破棄
terraform destroy
```

## 環境変数

### フロントエンド (.env.local)

```
NEXT_PUBLIC_API_URL=https://api.example.com
NEXT_PUBLIC_MAP_TILE_URL=https://cyberjapandata.gsi.go.jp/xyz/std/{z}/{x}/{y}.png
```

### バックエンド (Lambda環境変数)

```
DB_HOST=xxx.rds.amazonaws.com
DB_PORT=5432
DB_NAME=safearea
DB_USER=xxx
DB_PASSWORD=xxx
GEOCODING_API_URL=https://msearch.gsi.go.jp/address-search/AddressSearch
JSHIS_API_URL=https://www.j-shis.bosai.go.jp/
```

## セキュリティ注意事項

- RDSはPrivate Subnetに配置し、外部から直接アクセス不可
- Security Groupで5432ポートはLambdaからのみ許可
- 機密情報（DB接続情報等）は環境変数またはSecrets Managerで管理
- terraform.tfvarsはgitignore対象

## 開発フロー

1. Issueを作成
2. ブランチを切る (`feature/xxx`)
3. 実装
4. PR作成
5. レビュー
6. main にマージ
7. 自動デプロイ（GitHub Actions）

## 参考ドキュメント

- [README.md](./README.md) - プロジェクト詳細
- [docs/architecture.drawio](./docs/architecture.drawio) - インフラ構成図
- [docs/api.md](./docs/api.md) - API仕様（作成予定）
