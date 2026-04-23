FROM golang:alpine as builder  

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o mailApp ./cmd/api

RUN chmod +x /app/mailApp

FROM alpine:latest

RUN mkdir /app

WORKDIR /app

COPY --from=builder /app/mailApp .
COPY --from=builder /app/templates ./templates

CMD ["./mailApp"]