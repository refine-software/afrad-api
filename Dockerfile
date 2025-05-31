# building stage
FROM golang:1.24.3-alpine3.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && \
  go install github.com/swaggo/swag/cmd/swag@latest

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main cmd/api/main.go

# runtime stage
FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache tzdata

COPY --from=builder /app/main .

ENV PORT=8080
ENV APP_ENV=prod

EXPOSE 8080

CMD ["./main"]
