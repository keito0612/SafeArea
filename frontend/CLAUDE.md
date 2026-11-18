# SafeArea Frontend - Claude Code 指示書

## プロジェクト概要

| 項目 | 内容 |
|------|------|
| アプリ名 | SafeArea Frontend |
| 概要 | 災害リスク可視化Webアプリ |
| フレームワーク | Next.js 14 (App Router) |
| 言語 | TypeScript |

## 技術スタック

| カテゴリ | 技術 |
|----------|------|
| フレームワーク | Next.js 14 (App Router) |
| 言語 | TypeScript |
| スタイリング | Tailwind CSS |
| 状態管理 | useState + Context |
| データフェッチ | TanStack Query |
| HTTPクライアント | fetch (built-in) |
| フォーム | React Hook Form |
| バリデーション | Zod |
| 地図 | Leaflet / React-Leaflet |
| ユニットテスト | Vitest |
| コンポーネントテスト | React Testing Library |
| E2Eテスト | Playwright |

## アーキテクチャ

### 採用パターン：フロントエンド向けレイヤードアーキテクチャ

バックエンドと設計思想を統一し、将来の機能拡張に対応しやすい構成。

```
┌─────────────────────────────────────────────────────────┐
│  UI層（app/, components/）                              │
│  - ページ、コンポーネント                                │
│  - ユーザー操作の受付、表示                              │
└─────────────────────┬───────────────────────────────────┘
                      │ Hook呼び出し
┌─────────────────────▼───────────────────────────────────┐
│  Hooks層（hooks/）                                      │
│  - 状態管理（useState, useReducer）                     │
│  - UIとServiceの橋渡し                                  │
│  - ローディング・エラー状態の管理                        │
└─────────────────────┬───────────────────────────────────┘
                      │ Service呼び出し
┌─────────────────────▼───────────────────────────────────┐
│  Services層（services/）                                │
│  - ビジネスロジック                                     │
│  - データ変換・加工                                     │
│  - バリデーション                                       │
└─────────────────────┬───────────────────────────────────┘
                      │ Repository呼び出し
┌─────────────────────▼───────────────────────────────────┐
│  Repositories層（repositories/）                        │
│  - API通信のみ                                         │
│  - fetch/axios のラップ                                │
│  - エラーハンドリング（通信レベル）                      │
└─────────────────────┬───────────────────────────────────┘
                      │ HTTP
┌─────────────────────▼───────────────────────────────────┐
│  Backend API                                            │
└─────────────────────────────────────────────────────────┘
```

### 各層の責務

| 層 | ディレクトリ | 責務 | 依存先 |
|----|-------------|------|--------|
| UI | `app/`, `components/` | 表示、ユーザー操作 | Hooks |
| Hooks | `hooks/` | 状態管理、UIとServiceの接続 | Services |
| Services | `services/` | ビジネスロジック、データ変換 | Repositories |
| Repositories | `repositories/` | API通信 | なし（外部API） |

### 依存ルール

- **上位層 → 下位層** のみ依存可能
- **下位層 → 上位層** への依存は禁止
- **同一層内** での依存は最小限に

```
✅ Component → Hook → Service → Repository
❌ Repository → Service（逆方向は禁止）
❌ Service → Hook（逆方向は禁止）
```

### フォルダ構成

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
│
├── components/               # UI層（コンポーネント）
│   ├── ui/                  # 共通UI（Button, Card, Input等）
│   ├── map/                 # 地図関連
│   ├── search/              # 検索関連
│   ├── risk/                # リスク表示関連
│   └── layout/              # レイアウト（Header, Footer）
│
├── hooks/                    # Hooks層
│   ├── use-search.ts        # 災害リスク検索（詳細タブ用）
│   └── use-disaster-history.ts # 災害履歴取得（災害履歴タブ用）
│
├── services/                 # Services層
│   ├── search.service.ts    # 検索関連ロジック
│   └── disaster.service.ts  # 災害履歴関連ロジック
│
├── repositories/             # Repositories層
│   ├── search.repository.ts # 検索API
│   └── disaster.repository.ts # 災害履歴API
│
├── types/                    # 型定義
│   ├── index.ts
│   ├── risk.ts
│   └── disaster.ts
│
├── lib/                      # ユーティリティ
│   ├── api-client.ts        # fetch wrapper
│   └── utils.ts
│
├── __tests__/                # テスト
│   ├── unit/                # ユニットテスト
│   ├── components/          # コンポーネントテスト
│   └── e2e/                 # E2Eテスト
│
└── public/
```

## 画面構成

### 画面一覧

| # | パス | 画面名 | 概要 |
|---|------|--------|------|
| 1 | `/` | LP | サービス紹介（初めての人向け） |
| 2 | `/search` | 検索ページ | 検索欄 + 説明 |
| 3 | `/search/results` | 検索結果 | リスク表示 + 地図 |
| 4 | `/terms` | 利用規約 | 法的事項 |
| 5 | `/privacy` | プライバシーポリシー | 個人情報の取り扱い |

### 画面フロー

```
LP（/）
  │
  │ 「検索する」ボタン
  ▼
検索ページ（/search）
  │
  │ 住所を入力して検索
  ▼
検索結果（/search/results?address=xxx）
```

### 各画面の内容

#### 1. LP（`/`）
- サービスの説明
- 特徴・強み
- 使い方の流れ
- 「検索してみる」ボタン → `/search` へ

#### 2. 検索ページ（`/search`）
- 検索フォーム（上部）
- サービスの説明（下部）
  - 5種類の災害リスクを表示
  - ハザードマップを重ねて表示
  - 過去の災害履歴も確認可能

#### 3. 検索結果ページ（`/search/results`）

タブ切り替え式のUI構成

**詳細タブ（デフォルト）**
- 総合評価（高/中/低）
- タグ一覧
- 各災害リスク詳細（5種類）
- ハザードマップ（地図 + レイヤー切り替え）

**過去の災害履歴タブ**
- 災害履歴一覧（過去50年分）
- 災害種別アイコン（地震/台風/洪水等）
- 被害内容の表示
- タブに件数バッジ表示

**API呼び出し**
- 詳細タブ: `GET /api/search?address=xxx`（検索時に取得）
- 災害履歴タブ: `GET /api/search/disasters?address=xxx`（タブ切り替え時に遅延取得）

### レンダリング戦略

| ページ | 戦略 | 理由 |
|--------|------|------|
| LP (`/`) | SSG | 静的コンテンツ |
| 検索 (`/search`) | SSG + CSR | 静的部分 + フォーム |
| 検索結果 (`/search/results`) | CSR | クエリパラメータに依存 |
| 利用規約・プライバシー | SSG | 静的コンテンツ |

### Server Components vs Client Components

```
Server Components（デフォルト）
    - 静的ページ（LP, 利用規約等）
    - 初期データ取得

Client Components（'use client'）
    - フォーム入力
    - 地図操作
    - Hooks使用箇所
```

## コーディング規約

### 命名規則

| 対象 | 規則 | 例 |
|------|------|-----|
| ファイル名 | kebab-case | `search-form.tsx` |
| コンポーネント | PascalCase | `SearchForm` |
| 関数 | camelCase | `handleSubmit` |
| 定数 | UPPER_SNAKE_CASE | `API_BASE_URL` |
| 型/インターフェース | PascalCase | `RiskResponse` |

### コンポーネント

- 関数コンポーネントを使用（クラスコンポーネント禁止）
- Props は明示的に型定義
- `any` 禁止

```tsx
// Good
interface SearchFormProps {
  onSubmit: (address: string) => void;
  isLoading?: boolean;
}

export const SearchForm: React.FC<SearchFormProps> = ({ onSubmit, isLoading = false }) => {
  return (
    <form>...</form>
  );
};

// Bad
export function SearchForm(props: any) {
  return <form>...</form>;
}
```

### Server Components vs Client Components

```tsx
// Server Component（デフォルト）
// データ取得、静的表示に使用
export default async function Page() {
  const data = await fetchData();
  return <div>{data}</div>;
}

// Client Component（インタラクションが必要な場合のみ）
'use client';

export function SearchForm() {
  const [value, setValue] = useState('');
  return <input value={value} onChange={(e) => setValue(e.target.value)} />;
}
```

### インポート順序

```tsx
// 1. React/Next.js
import { useState } from 'react';
import Link from 'next/link';

// 2. 外部ライブラリ
import { useForm } from 'react-hook-form';

// 3. 内部モジュール（絶対パス）
import { Button } from '@/components/ui/button';
import { fetchRisk } from '@/lib/api';

// 4. 型
import type { RiskResponse } from '@/types';
```

### Tailwind CSS

- インラインで記述（CSS Modules不使用）
- 長いクラス名は改行して整理

```tsx
// Good
<button
  className={cn(
    'px-4 py-2 rounded-lg font-medium',
    'bg-blue-600 text-white',
    'hover:bg-blue-700 transition-colors',
    isLoading && 'opacity-50 cursor-not-allowed'
  )}
>
  検索
</button>

// Bad - 1行に詰め込む
<button className="px-4 py-2 rounded-lg font-medium bg-blue-600 text-white hover:bg-blue-700 transition-colors">
```

## API クライアント

### 基本構造

```typescript
// lib/api.ts
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL;

export async function fetchRisk(address: string): Promise<RiskResponse> {
  const res = await fetch(`${API_BASE_URL}/api/risk?address=${encodeURIComponent(address)}`);

  if (!res.ok) {
    throw new Error('Failed to fetch risk data');
  }

  return res.json();
}
```

### エラーハンドリング

```typescript
// コンポーネント側
try {
  const data = await fetchRisk(address);
  setResult(data);
} catch (error) {
  setError('災害リスクの取得に失敗しました');
}
```

## 型定義

### APIレスポンス

```typescript
// types/index.ts
export interface RiskResponse {
  success: boolean;
  data: {
    overall_risk: 'high' | 'medium' | 'low';
    risks: Risk[];
    tags: string[];
    disasters: Disaster[];
  };
}

export interface Risk {
  type: 'flood' | 'landslide' | 'earthquake' | 'tsunami' | 'liquefaction';
  level: 'high' | 'medium' | 'low' | 'none';
  description: string;
}

export interface Disaster {
  id: string;
  name: string;
  year: number;
  description: string;
  article_url?: string;
}
```

## 地図（Leaflet）

### 注意点

- `react-leaflet` は Client Component でのみ使用可能
- SSR を避けるため `dynamic import` を使用

```tsx
// components/map/hazard-map.tsx
'use client';

import dynamic from 'next/dynamic';

const MapContainer = dynamic(
  () => import('react-leaflet').then((mod) => mod.MapContainer),
  { ssr: false }
);
```

### ハザードマップタイル

```typescript
// 国土地理院のハザードマップタイル
const HAZARD_TILES = {
  flood: 'https://disaportaldata.gsi.go.jp/raster/01_flood_l2_shinsuishin_data/{z}/{x}/{y}.png',
  landslide: 'https://disaportaldata.gsi.go.jp/raster/05_dosekiryukeikaikuiki/{z}/{x}/{y}.png',
  tsunami: 'https://disaportaldata.gsi.go.jp/raster/04_tsunami_newlegend_data/{z}/{x}/{y}.png',
};
```

## 開発コマンド

```bash
# 開発サーバー
npm run dev

# ビルド
npm run build

# Lint
npm run lint

# 型チェック
npm run type-check

# フォーマット
npm run format
```

## 環境変数

| 変数名 | 説明 |
|--------|------|
| `NEXT_PUBLIC_API_URL` | バックエンドAPI URL |
| `NEXT_PUBLIC_MAP_TILE_URL` | 地図タイルURL |

## テスト

### テストライブラリ

| 種類 | ライブラリ | 用途 |
|------|-----------|------|
| ユニット | Vitest | Services, Repositories, ユーティリティ |
| コンポーネント | React Testing Library | コンポーネントの振る舞い |
| E2E | Playwright | 画面遷移、ユーザーフロー |

### テストコマンド

```bash
# ユニット/コンポーネントテスト
npm run test

# ウォッチモード
npm run test:watch

# カバレッジ
npm run test:coverage

# E2Eテスト
npm run test:e2e

# E2Eテスト（UIモード）
npm run test:e2e:ui
```

### テストファイルの配置

```
__tests__/
├── unit/                    # ユニットテスト
│   ├── services/
│   └── repositories/
├── components/              # コンポーネントテスト
│   ├── ui/
│   └── search/
└── e2e/                     # E2Eテスト
    ├── search.spec.ts
    └── navigation.spec.ts
```

## 参考ドキュメント

- [README.md](./README.md) - セットアップ手順
- [../README.md](../README.md) - プロジェクト全体の仕様
