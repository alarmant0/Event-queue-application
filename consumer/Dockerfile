FROM golang:1.17-alpine AS builder

ENV GO111MODULE=on \
    GOOS=linux \
    CGO_ENABLE=0 \
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
