FROM golang:1.17

WORKDIR /api

COPY go.mod .
COPY go.sum	.
RUN go mod download

COPY . .

RUN go get github.com/codegangsta/gin

CMD ["gin","run","main.go"]
