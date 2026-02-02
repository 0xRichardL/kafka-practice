FROM golang:1.25.4 as builder

WORKDIR /app

COPY . .

RUN go build -o main .

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache libc6-compat

COPY --from=builder /app/main .

CMD ["./main"]