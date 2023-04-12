FROM golang:alpine

RUN apk update && apk add --no-cache git && apk add --no-cache bash && apk add build-base

WORKDIR /app
COPY go.mod go.sum ./
COPY . .

RUN go mod download
RUN go build -o /api ./cmd/api

EXPOSE 8080
CMD ["/api"]