FROM golang:alpine AS builder

WORKDIR /app

# mattn/go-sqlite3 uses CGO, so the builder needs a C toolchain.
RUN apk --no-cache add build-base

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 go build -o scraper-api ./cmd/main.go

FROM alpine:latest

# BusyBox provides crond; sqlite is used by entrypoint.sh for the first-run count.
RUN apk --no-cache add ca-certificates sqlite

WORKDIR /root/

COPY --from=builder /app/scraper-api .
COPY config.json .
COPY entrypoint.sh .

RUN chmod +x entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["./entrypoint.sh"]
