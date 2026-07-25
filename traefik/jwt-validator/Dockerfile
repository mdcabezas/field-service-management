FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY main.go go.mod ./
RUN CGO_ENABLED=0 GOOS=linux go build -o jwt-validator main.go

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/jwt-validator /usr/local/bin/
EXPOSE 8080
CMD ["jwt-validator"]
