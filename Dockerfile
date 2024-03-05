FROM golang:1.22.1

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY /cmd* ./cmd/
COPY /internal* ./internal/

RUN go build -o legbot ./cmd/legbot.go

CMD ["./legbot"]