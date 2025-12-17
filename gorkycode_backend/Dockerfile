# syntax=docker/dockerfile:1

FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . .

RUN go build -o gorkycode-backend ./cmd/main/main.go

FROM alpine:3.18

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/gorkycode-backend .

EXPOSE 8080

CMD ["./gorkycode-backend"]