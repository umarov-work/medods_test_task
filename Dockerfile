FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o medods_test_task ./cmd

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/medods_test_task .

RUN apt-get update && apt-get install -y ca-certificates

EXPOSE 8080

CMD ["./medods_test_task"]
