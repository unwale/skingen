FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o task-service ./services/task-service/cmd/task-service

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/task-service .

EXPOSE ${TASK_SERVICE_PORT}

CMD ["./task-service"]