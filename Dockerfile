# ビルドステージ
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main ./cmd/api
RUN go build -o migrate ./cmd/migrate

# 実行ステージ
FROM alpine:3.24
WORKDIR /app
# ビルドステージから実行ファイルをコピー
COPY --from=builder /app/main ./main
COPY --from=builder /app/migrate ./migrate
COPY --from=builder /app/sql ./sql
# ポートを公開
EXPOSE 8080
CMD ["./main"]
