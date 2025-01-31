FROM golang:1.23.4 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/google/wire/cmd/wire@latest

COPY . .

RUN wire ./internal/server

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/payment


FROM alpine:latest

WORKDIR /root/


RUN apk add --no-cache go

COPY --from=builder /app/main .

COPY --from=builder /app/config/.env ./config/.env

CMD ["./main"]
