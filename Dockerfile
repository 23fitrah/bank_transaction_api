FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod tidy
RUN GOOS=linux GOARCH=amd64 go build -o transaction_api .

FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/transaction_api .
COPY --from=builder /app/.env .

RUN chmod +x ./transaction_api

EXPOSE 6969

CMD ["./transaction_api"]
