FROM golang:1.19

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o main ./cmd

EXPOSE 5672

EXPOSE 9092

EXPOSE 15672

CMD ["./main"]