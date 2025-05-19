# ───── Stage 1: Build ─────
FROM golang:1.24.1-bullseye AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO should be enabled
RUN CGO_ENABLED=1 go build -o /go/bin/server ./cmd/application

# ───── Stage 2: Runtime ─────
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y \
    ca-certificates \
    libsqlite3-0 \
 && rm -rf /var/lib/apt/lists/*

WORKDIR /app
RUN mkdir -p data

COPY --from=builder /go/bin/server /app/server

EXPOSE 8080

# run server
CMD ["/app/server"]
