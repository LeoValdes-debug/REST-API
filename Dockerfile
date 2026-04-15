# Многоэтапная сборка (multi-stage build).
# Этап 1 — сборка бинаря. Этап 2 — минимальный образ только с бинарём.
# Итоговый образ ~10MB вместо ~800MB если бы оставили весь Go-toolchain.

FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
