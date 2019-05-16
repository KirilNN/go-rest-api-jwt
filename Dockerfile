FROM golang:1.12.5-alpine3.9

RUN mkdir /app

COPY main.go /app

WORKDIR /app

RUN go build -o main .

CMD ["/app/main"]