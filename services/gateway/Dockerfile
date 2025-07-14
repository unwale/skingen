FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o gateway ./services/gateway/cmd/gateway

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/gateway .

EXPOSE ${GATEWAY_PORT}

CMD ["./gateway"]