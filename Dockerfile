# ビルドステージ
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main ./cmd/api/main.go

# 実行ステージ
FROM alpine:3.20
WORKDIR /app
# ビルドステージから実行ファイルをコピー
COPY --from=builder /app/main ./main
# ポートを公開
EXPOSE 8080
CMD ["./main"]
