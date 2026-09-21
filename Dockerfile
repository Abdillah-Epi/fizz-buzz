FROM golang:1.26.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o fizzbuzz-api ./cmd/api

FROM debian:12-slim AS server

COPY --from=builder /app/fizzbuzz-api /fizzbuzz-api

EXPOSE 8080

CMD ["/fizzbuzz-api"]
