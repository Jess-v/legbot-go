FROM golang:1.19.0

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY /cmd* ./cmd/
COPY /internal* ./internal/

RUN go build -o legbot ./cmd/legbot.go

CMD ["./legbot"]