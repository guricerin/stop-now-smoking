FROM golang:1.17.1-alpine

WORKDIR /go/src/app

RUN apk update \
    && apk add git

COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

CMD [ "go", "run", "./cmd/local/main.go" ]
