FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o backend ./cmd/server

FROM alpine:3.21
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/backend /usr/local/bin/
EXPOSE 8080
CMD ["backend"]
