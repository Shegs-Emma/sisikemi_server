# Build stage
FROM golang:1.23-alpine3.20 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.20
WORKDIR /app

RUN apk add --no-cache curl ca-certificates postgresql-client bash

RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz \
  | tar -xz && \
  mv migrate /usr/local/bin/

COPY --from=builder /app/main .
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./db/migration

RUN chmod +x start.sh wait-for.sh

EXPOSE 8080
ENTRYPOINT ["/app/start.sh"]
