FROM golang:1.26.5-alpine3.24 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o mock-backend ./deploy/mock-backend

FROM alpine:3.24
COPY --from=builder /app/mock-backend /mock-backend
EXPOSE 8081
CMD ["/mock-backend"]
