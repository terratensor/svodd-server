# syntax = docker/dockerfile:1.2
FROM golang:latest AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o ./svodd-server-parser ./cmd/parser

#2

FROM scratch

WORKDIR /app

COPY --from=builder /app/svodd-server-parser /app/svodd-server-parser
COPY --from=builder /app/config /app/config
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Корневой сертификат удостоверяющего центра (УЦ) Минцифры
COPY --from=builder /app/certs/ /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
ENV TZ=Europe/Moscow

CMD ["./svodd-server-parser"]