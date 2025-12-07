FROM golang:1.25.5-alpine as builder

WORKDIR /app/producer

COPY . .

RUN go build -o main 

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 5001

CMD ["./main"]