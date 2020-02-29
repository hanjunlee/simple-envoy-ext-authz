FROM golang:1.13

WORKDIR /app

COPY go.mod go.sum /app/

RUN go mod download

COPY . /app/

RUN go build -o ext-authz cmd/main.go

ENTRYPOINT ["./ext-authz", "server"]
