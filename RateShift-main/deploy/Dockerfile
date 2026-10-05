FROM golang:1.26.5-alpine3.24 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o gateway ./cmd/gateway

FROM alpine:3.24
COPY --from=builder /app/gateway /gateway
EXPOSE 8080
CMD ["/gateway"]
