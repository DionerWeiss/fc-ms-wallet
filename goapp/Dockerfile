FROM golang:1.20 as builder

WORKDIR /app

RUN go mod init go-challenge
COPY . .
RUN go build  -o /app/main ./cmd/walletcore/


FROM scratch
WORKDIR /app
COPY --from=builder /app .
CMD ["/app/main"]

# WORKDIR /app

# RUN apt-get update && apt-get install -y librdkafka-dev

# # RUN go mod init goapp

# # COPY go.mod ./
# # COPY go.sum ./
# COPY . .

# RUN go mod download

# # RUN cd cmd/walletcore

# COPY *.go ./

# RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/walletcore


# FROM scratch
# WORKDIR / 
# COPY --from=builder /app/main /main
# ENTRYPOINT ["./cmd/walletcore/main"]