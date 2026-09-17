FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY main-service/go.mod main-service/go.sum ./
RUN go mod download

COPY main-service/ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/main-service ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates wget

WORKDIR /app

COPY --from=builder /app/bin/main-service .

EXPOSE 8080

CMD ["./main-service"]