FROM golang:1.17.6-alpine AS builder

WORKDIR /app
COPY . .

RUN apk --no-cache add ca-certificates

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "main" -ldflags="-w -s" ./main.go

CMD ["/app/main"]