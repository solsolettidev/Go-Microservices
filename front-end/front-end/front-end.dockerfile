FROM golang:1.18-alpine as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o frontApp ./cmd/web

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/frontApp .
COPY --from=builder /app/cmd/web/templates ./cmd/web/templates

CMD ["./frontApp"]
