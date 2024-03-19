FROM golang:1.17-alpine AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GO111MODULE=on \
    GOARCH=amd64

WORKDIR /app

COPY ./go.mod ./go.sum ./

RUN go mod download

COPY ./main.go ./

RUN go build -o main .

FROM alpine:latest

COPY --from=builder /app/main /app/main

EXPOSE 8080

CMD ["/app/main"]
