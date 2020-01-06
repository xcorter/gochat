FROM golang:1.13-alpine

WORKDIR /app

COPY ./ /app

RUN go mod download

RUN go get github.com/githubnemo/CompileDaemon

ENTRYPOINT go run cmd/chat/main.go
#ENTRYPOINT CompileDaemon --build="go build ./cmd/chat/main.go" --command=./main