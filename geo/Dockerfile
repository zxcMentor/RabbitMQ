FROM golang:1.19

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o main ./cmd

EXPOSE 8081

EXPOSE 5672

EXPOSE 50051

CMD ["./main"]