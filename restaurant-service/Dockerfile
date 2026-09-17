FROM golang:1.21-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum
COPY restaurant-service/go.mod restaurant-service/go.sum ./
RUN go mod download

# Копируем исходный код
COPY restaurant-service/ ./

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/restaurant-service ./cmd/main.go

# Финальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/bin/restaurant-service .

EXPOSE 8081

CMD ["./restaurant-service"]