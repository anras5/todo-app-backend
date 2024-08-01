FROM golang:alpine

RUN apk update && apk add --no-cache git && apk add --no-cache bash && apk add build-base

WORKDIR /app
COPY go.mod go.sum ./
COPY . .

RUN go mod download
RUN go install github.com/air-verse/air@latest


CMD air --build.cmd "go build -o apibin ./cmd/api" --build.bin "./apibin" --build.exclude_dir "postgres-data"
