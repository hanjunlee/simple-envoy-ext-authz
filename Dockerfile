FROM golang:1.13

WORKDIR /app

COPY go.mod go.sum /app/

RUN go mod download

COPY . /app/

RUN go build -o server main.go

ENTRYPOINT ["./server"]
