FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o analytics ./cmd/server

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/analytics /usr/local/bin/analytics
EXPOSE 8081
ENTRYPOINT ["/usr/local/bin/analytics"]
