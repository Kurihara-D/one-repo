FROM golang:1.21-bullseye AS build

WORKDIR /go/src/app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .
EXPOSE 3000

ENV GO111MODULE=on

RUN go get -u github.com/cosmtrek/air && \
    go build -o /go/bin/air github.com/cosmtrek/air
#CMD ["go", "run", "main.go"]
CMD ["air", "-c", ".air.toml"]