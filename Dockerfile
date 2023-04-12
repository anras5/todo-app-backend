FROM golang:alpine

RUN apk update && apk add --no-cache git && apk add --no-cache bash && apk add build-base

WORKDIR /app
COPY go.mod go.sum ./
COPY . .

RUN go mod download
RUN go get github.com/githubnemo/CompileDaemon
RUN go install github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon -polling -log-prefix=false -build="go build -o api ./cmd/api" -command="./api" -directory="./"