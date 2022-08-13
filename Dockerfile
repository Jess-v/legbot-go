FROM golang:1.19.0

WORKDIR /app

COPY . .

RUN go get -d -v 

RUN go install -v

RUN mkdir users

CMD ["main"]