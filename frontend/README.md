# SafeArea Frontend

災害リスク可視化サービス「SafeArea」のフロントエンド

## 概要

| 項目 | 内容 |
|------|------|
| フレームワーク | Next.js 14 (App Router) |
| 言語 | TypeScript |
| アーキテクチャ | レイヤードアーキテクチャ |

## 技術スタック

| カテゴリ | 技術 |
|----------|------|
| スタイリング | Tailwind CSS |
| 状態管理 | useState + Context |
| データフェッチ | TanStack Query |
| フォーム | React Hook Form |
| バリデーション | Zod |
| 地図 | Leaflet / React-Leaflet |
| テスト | Vitest / React Testing Library / Playwright |

## セットアップ

### 必要なもの

- Node.js 20+
- npm or pnpm

### ローカル環境構築

```bash
# リポジトリをクローン（モノレポ）
git clone https://github.com/isobekeito/safearea.git
cd safearea

# 全サービスを起動（推奨）
docker-compose up -d
# http://localhost:3000 でアクセス

# または、フロントエンドのみ個別起動
cd frontend
npm install
npm run dev
```

### 環境変数

`.env.local` を作成:

```
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_MAP_TILE_URL=https://cyberjapandata.gsi.go.jp/xyz/std/{z}/{x}/{y}.png
```

## 開発コマンド

```bash
# 開発サーバー
npm run dev

# ビルド
npm run build

# 本番起動
npm run start

# Lint
npm run lint

# 型チェック
npm run type-check

# テスト
npm run test

# E2Eテスト
npm run test:e2e
```

## ディレクトリ構成

```
frontend/
├── app/                      # UI層（ページ）
│   ├── layout.tsx
│   ├── page.tsx             # LP
│   ├── search/
│   │   ├── page.tsx         # 検索ページ
│   │   └── results/
│   │       └── page.tsx     # 検索結果ページ
│   ├── terms/
│   │   └── page.tsx         # 利用規約
│   └── privacy/
│       └── page.tsx         # プライバシーポリシー
├── components/               # UI層（コンポーネント）
│   ├── ui/                  # 共通UI
│   ├── map/                 # 地図関連
│   ├── search/              # 検索関連
│   ├── risk/                # リスク表示関連
│   └── layout/              # レイアウト
├── hooks/                    # Hooks層
├── services/                 # Services層
├── repositories/             # Repositories層
├── types/                    # 型定義
├── lib/                      # ユーティリティ
├── __tests__/                # テスト
│   ├── unit/
│   ├── components/
│   └── e2e/
└── public/
```

## 画面構成

| # | パス | 画面 | 説明 |
|---|------|------|------|
| 1 | `/` | LP | サービス紹介（初めての人向け） |
| 2 | `/search` | 検索ページ | 検索欄 + 説明 |
| 3 | `/search/results` | 検索結果 | 災害リスク表示 + 地図 |
| 4 | `/terms` | 利用規約 | 法的事項 |
| 5 | `/privacy` | プライバシーポリシー | 個人情報の取り扱い |

## 画面フロー

```
LP（/）→ 検索ページ（/search）→ 検索結果（/search/results?address=xxx）
```

## 機能一覧

### LP（`/`）

| # | 機能 | 説明 | 優先度 |
|---|------|------|:------:|
| 1-1 | サービス説明表示 | SafeAreaの概要を表示 | MVP |
| 1-2 | 特徴紹介 | 3つの強みを表示 | MVP |
| 1-3 | 使い方説明 | 3ステップの使い方を表示 | MVP |
| 1-4 | 対応災害一覧 | 5種類の災害リスクを紹介 | MVP |
| 1-5 | 検索ページへの導線 | 「検索してみる」ボタン | MVP |

### 検索ページ（`/search`）

| # | 機能 | 説明 | 優先度 |
|---|------|------|:------:|
| 2-1 | 住所入力フォーム | テキスト入力で住所を入力 | MVP |
| 2-2 | 検索ボタン | 入力した住所でリスク検索 | MVP |
| 2-3 | 入力バリデーション | 空欄チェック、文字数チェック | MVP |
| 2-4 | サービス説明表示 | 検索フォーム下に説明を表示 | MVP |
| 2-5 | 住所候補サジェスト | 入力中に候補を表示 | 将来 |

### 検索結果ページ（`/search/results`）

タブ切り替え式のUI。詳細タブと過去の災害履歴タブで構成。

#### 共通

| # | 機能 | 説明 | 優先度 |
|---|------|------|:------:|
| 3-0 | タブ切り替え | 詳細/災害履歴タブの切り替え | MVP |
| 3-1 | 再検索 | 別の住所で再検索 | MVP |

#### 詳細タブ

| # | 機能 | 説明 | 優先度 |
|---|------|------|:------:|
| 3-2 | 総合評価表示 | 高/中/低の3段階で表示 | MVP |
| 3-3 | タグ表示 | 該当するリスクをタグで表示 | MVP |
| 3-4 | 災害リスク詳細表示 | 5種類のリスクを個別表示 | MVP |
| 3-5 | ハザードマップ表示 | 地図上にリスクを重ねて表示 | MVP |
| 3-6 | 地図レイヤー切り替え | 洪水/土砂/津波の切り替え | MVP |

#### 過去の災害履歴タブ

| # | 機能 | 説明 | 優先度 |
|---|------|------|:------:|
| 3-7 | 災害履歴一覧表示 | 過去50年の災害履歴を一覧表示 | MVP |
| 3-8 | 災害種別アイコン | 地震/台風/洪水等をアイコンで表示 | MVP |
| 3-9 | 災害履歴件数バッジ | タブに件数を表示 | MVP |

#### 将来機能

| # | 機能 | 説明 | 優先度 |
|---|------|------|:------:|
| 3-10 | 結果シェア | SNSシェアボタン | 将来 |
| 3-11 | PDF出力 | 結果をPDFで出力 | 将来 |

### 共通機能

| # | 機能 | 説明 | 優先度 |
|---|------|------|:------:|
| C-1 | ヘッダー | ロゴ、ナビゲーション | MVP |
| C-2 | フッター | リンク、コピーライト | MVP |
| C-3 | レスポンシブ対応 | モバイル/タブレット/PC対応 | MVP |
| C-4 | ローディング表示 | API通信中の表示 | MVP |
| C-5 | エラー表示 | エラー時のメッセージ表示 | MVP |
| C-6 | 404ページ | 存在しないページの表示 | MVP |

### 災害リスク種別

| # | 種別 | アイコン | 説明 |
|---|------|:--------:|------|
| R-1 | 洪水浸水 | 🌊 | 河川氾濫による浸水リスク |
| R-2 | 土砂災害 | 🏔️ | 土砂災害警戒区域 |
| R-3 | 地震危険度 | 🌏 | 地震による被害リスク |
| R-4 | 津波浸水 | 🌊 | 津波による浸水リスク |
| R-5 | 液状化 | 💧 | 地盤の液状化リスク |

## API連携

バックエンドAPI（Go Lambda）と連携:

| エンドポイント | 用途 | 呼び出しタイミング |
|---------------|------|-------------------|
| `GET /api/search?address=xxx` | 災害リスク検索 | 検索実行時 |
| `GET /api/search/disasters?address=xxx` | 過去の災害履歴 | 災害履歴タブ切り替え時 |

## デプロイ

AWS Amplify Hosting にデプロイ

```bash
# ビルド
npm run build

# Amplifyが自動でデプロイ（GitHub連携）
```

## 関連リポジトリ

- [backend/](../backend/) - バックエンドAPI（Go）
- [infrastructure/](../infrastructure/) - インフラ（Terraform）
