# Builder
FROM golang:1.22 AS builder

WORKDIR /usr/local/src
COPY go.mod .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./.bin/warehouse-api ./cmd/warehouse-api/main.go

# App runner
FROM alpine:latest

WORKDIR /usr/local/src

COPY --from=builder /usr/local/src/.bin/warehouse-api /usr/local/src/.bin/warehouse-api
COPY --from=builder /usr/local/src/.env /usr/local/src/
COPY --from=builder /usr/local/src/configs/main.yml /usr/local/src/configs/
COPY --from=builder /usr/local/src/migrations /usr/local/src/migrations/

CMD ["./.bin/warehouse-api"]