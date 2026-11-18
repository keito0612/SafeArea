# SafeArea Backend

災害リスク可視化サービス「SafeArea」のバックエンドAPI

## 概要

| 項目 | 内容 |
|------|------|
| 言語 | Go 1.22 |
| ランタイム | AWS Lambda |
| ルーター | chi |
| ORM | GORM |
| DB | PostgreSQL |

## セットアップ

### 必要なもの

- Go 1.22+
- Docker / Docker Compose
- make (任意)

### ローカル環境構築

```bash
# リポジトリをクローン（モノレポ）
git clone https://github.com/isobekeito/safearea.git
cd safearea

# 全サービスを起動（推奨）
docker-compose up -d

# または、バックエンドのみ個別起動
cd backend
go mod download
docker-compose up -d  # DBのみ起動
go run cmd/api/main.go
```

### 環境変数

`.env.local` を作成:

```
ENV=local
DB_HOST=localhost
DB_PORT=5432
DB_NAME=safearea
DB_USER=postgres
DB_PASSWORD=postgres
GEOCODING_API_URL=https://msearch.gsi.go.jp/address-search/AddressSearch
JSHIS_API_URL=https://www.j-shis.bosai.go.jp/
```

## 開発コマンド

```bash
# ローカル実行
go run cmd/api/main.go

# テスト
go test ./...

# テスト（カバレッジ付き）
go test -cover ./...

# Lint
golangci-lint run

# ビルド（Lambda用）
GOOS=linux GOARCH=amd64 go build -o bootstrap cmd/api/main.go
```

## API エンドポイント

| Method | Path | 説明 |
|--------|------|------|
| GET | `/api/health` | ヘルスチェック |
| GET | `/api/search?address=xxx` | 災害リスク検索（詳細タブ用） |
| GET | `/api/search/disasters?address=xxx` | 過去の災害履歴（災害履歴タブ用） |

## 機能一覧

| # | 機能 | エンドポイント | 説明 | 優先度 |
|---|------|---------------|------|:------:|
| B-1 | ヘルスチェック | GET `/api/health` | サーバー稼働確認 | MVP |
| B-2 | 災害リスク検索 | GET `/api/search` | 住所から5種類のリスクを取得 | MVP |
| B-3 | 過去の災害履歴 | GET `/api/search/disasters` | 住所周辺の災害履歴を取得 | MVP |
| B-4 | ジオコーディング | 内部処理 | 住所→緯度経度変換（国土地理院API） | MVP |
| B-5 | リスク判定 | 内部処理 | 外部APIからリスクレベルを算出 | MVP |

## APIレスポンス

### GET `/api/search?address=xxx`

詳細タブ用のリスク情報を返却

```json
{
  "success": true,
  "data": {
    "address": "東京都江東区豊洲1丁目",
    "coordinates": { "lat": 35.655, "lng": 139.793 },
    "overall_risk": "high",
    "risks": [
      { "type": "flood", "level": "high", "description": "浸水想定区域内" },
      { "type": "landslide", "level": "low", "description": "警戒区域外" },
      { "type": "earthquake", "level": "medium", "description": "震度6弱の可能性" },
      { "type": "tsunami", "level": "medium", "description": "浸水想定1m未満" },
      { "type": "liquefaction", "level": "high", "description": "液状化の可能性が高い" }
    ],
    "tags": ["flood", "liquefaction"],
    "disaster_history_count": 3
  }
}
```

### GET `/api/search/disasters?address=xxx`

災害履歴タブ用の過去の災害一覧を返却（過去50年分）

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
      },
      {
        "id": "d002",
        "type": "typhoon",
        "name": "台風19号",
        "year": 2019,
        "damage": "高潮警戒、一部地域で浸水"
      }
    ]
  }
}
```

### 災害種別（type）

| type | 説明 |
|------|------|
| earthquake | 地震 |
| typhoon | 台風 |
| flood | 洪水 |
| landslide | 土砂災害 |
| tsunami | 津波 |
| other | その他 |

## アーキテクチャ

詳細は [ARCHITECTURE.md](./ARCHITECTURE.md) を参照

```
internal/
├── controller/       # コントローラー層
├── usecase/          # ユースケース層
├── domain/           # ドメイン層
│   ├── model/
│   └── repository/
├── infrastructure/   # インフラストラクチャ層
│   ├── persistence/
│   └── external/
└── dto/              # DTO
```

## ディレクトリ構成

```
backend/
├── cmd/api/              # エントリーポイント
├── internal/             # アプリケーションコード
├── db/                   # マイグレーション
├── docker-compose.yml    # ローカルDB
├── go.mod
├── ARCHITECTURE.md       # アーキテクチャ詳細
├── CLAUDE.md             # Claude Code指示書
└── README.md             # このファイル
```

## デプロイ

AWS Lambda + API Gateway にデプロイ

```bash
# ビルド
GOOS=linux GOARCH=amd64 go build -o bootstrap cmd/api/main.go

# ZIPに圧縮
zip deployment.zip bootstrap

# Terraformでデプロイ（infrastructureリポジトリ）
```

## 災害履歴データの取得・登録

### データ取得フロー

```
1. CSVダウンロード
   内閣府・気象庁のサイトからCSV/Excelをダウンロード
   ↓
2. データ加工（スクリプト）
   - CSVをパース
   - 必要なカラムを抽出
   - 指定形式に変換
   - 過去50年分（1975年〜）にフィルタリング
   ↓
3. DBに保存
   シードスクリプトで投入
```

### データソース

| ソース | URL | 形式 |
|--------|-----|------|
| 内閣府 災害関係データ | https://www.bousai.go.jp/shiryou/data.html | PDF/Excel |
| 気象庁 震度データベース | https://www.data.jma.go.jp/svd/eqdb/data/shindo/ | CSV |
| 気象庁 災害をもたらした気象事例 | https://www.data.jma.go.jp/obd/stats/data/bosai/report/ | PDF/CSV |

### 関連ファイル

```
backend/
├── db/
│   ├── migrations/           # テーブル作成
│   └── seeds/                # 初期データ
│       └── disasters.csv     # 加工済み災害データ
├── scripts/
│   └── import-disasters/     # CSV加工スクリプト
│       ├── main.go
│       └── README.md
```

### CSV加工スクリプトの処理

```go
// 1. 元CSVを読み込み
// 2. 必要なカラムを抽出（災害名、発生年、都道府県、被害内容）
// 3. 災害種別を判定
//    - "地震" → earthquake
//    - "台風" → typhoon
//    - "豪雨" "洪水" → flood
//    - "土砂" → landslide
//    - "津波" → tsunami
// 4. 過去50年分にフィルタリング（year >= 1975）
// 5. 指定形式で出力
```

### 加工後のCSV形式

```csv
id,type,name,year,damage,prefecture_code,city_name
uuid1,earthquake,東日本大震災,2011,液状化により道路陥没,13,江東区
uuid2,typhoon,台風19号,2019,高潮警戒・一部浸水,13,江東区
```

### シード実行

```bash
# 災害データをDBに投入
go run scripts/import-disasters/main.go

# または make コマンド
make seed-disasters
```

## 関連ディレクトリ

- [../frontend/](../frontend/) - フロントエンド（Next.js）
- [../infrastructure/](../infrastructure/) - インフラ（Terraform）※作成予定
