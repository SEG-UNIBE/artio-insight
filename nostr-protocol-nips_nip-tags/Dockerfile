FROM golang:latest

WORKDIR /app
COPY ./nostr-protocol-nips_nip-tags /app
COPY .env /app/.env

RUN go mod tidy
RUN go mod download
RUN go build -o /app/bin/main /app/cmd/main.go


CMD ["/app/bin/main"]