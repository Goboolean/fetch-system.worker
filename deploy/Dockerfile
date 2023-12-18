FROM golang

WORKDIR /go/src/app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o . cmd/app/main.go

CMD ["./main"]