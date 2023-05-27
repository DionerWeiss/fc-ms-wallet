FROM golang:1.20

WORKDIR /app/

RUN apt-get update && apt-get install -y librdkafka-dev
COPY . .
RUN go mod tidy

EXPOSE 8080

CMD ["go", "run", "cmd/walletcore/main.go"]