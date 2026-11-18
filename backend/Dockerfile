# Backend Dockerfile (Go)
FROM golang:1.23-alpine

WORKDIR /app

# 必要なツールをインストール
RUN apk add --no-cache git

# Air（ホットリロード）をインストール
RUN go install github.com/air-verse/air@v1.52.3

# go.modをコピーして依存関係をダウンロード
COPY go.mod ./
RUN go mod download || true

# ソースコードをコピー
COPY . .

# go.sumを生成（存在しない場合）
RUN go mod tidy

# ポートを公開
EXPOSE 8080

# 開発時はAirでホットリロード
CMD ["air", "-c", ".air.toml"]
