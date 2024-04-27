FROM golang:1.22 AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o build/ cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/build/main .

RUN chmod +x main

EXPOSE 8080

CMD ["./main"]