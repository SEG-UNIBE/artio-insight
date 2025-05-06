FROM golang:latest

WORKDIR /app
COPY . /app

RUN go mod tidy
RUN go mod download
RUN go build -o /app/bin/go-db-etl /app/cmd/main.go


CMD ["/app/bin/go-db-etl"]