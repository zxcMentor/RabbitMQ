FROM golang:1.19-alpine as builder

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o main ./cmd

FROM alpine

COPY --from=builder /app/main /main

EXPOSE 8082
EXPOSE 50052

CMD ["/main"]