FROM golang:alpine3.18

WORKDIR /app

# Pra usar o alpine, precisamos desta dependência.
RUN apk add alpine-sdk

COPY . .
RUN go mod tidy
# A documentação o Kafka pede para passar o parâmetro -tags musl pra o build funcionar neste SO.
RUN go build -tags musl -o /app/main ./cmd/walletcore/

CMD ["/app/main"]

# FROM golang:1.20 as builder

# WORKDIR /app

# RUN go mod init go-challenge
# COPY . .
# RUN go build  -o main ./cmd/walletcore/


# FROM scratch
# WORKDIR /app
# COPY --from=builder /app .
# CMD ["main"]

# # WORKDIR /app

# # RUN apt-get update && apt-get install -y librdkafka-dev

# # # RUN go mod init goapp

# # # COPY go.mod ./
# # # COPY go.sum ./
# # COPY . .

# # RUN go mod download

# # # RUN cd cmd/walletcore

# # COPY *.go ./

# # RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/walletcore


# # FROM scratch
# # WORKDIR / 
# # COPY --from=builder /app/main /main
# # ENTRYPOINT ["./cmd/walletcore/main"]