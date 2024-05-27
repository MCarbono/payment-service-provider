# Build stage
FROM golang:1.19-alpine3.16 AS builder
WORKDIR /app
RUN apk add --no-cache curl
COPY . .
RUN go build -o main main.go
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.12.2/migrate.linux-amd64.tar.gz | tar xvz

#Run stage
FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate.linux-amd64 ./migrate
COPY .env .
COPY infra/db/migration ./infra/db/migration

EXPOSE 3000
ENTRYPOINT ["/app/main", "--env=production"]