FROM golang

WORKDIR /go/src/app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o . cmd/main/run.go

CMD ["./run"]