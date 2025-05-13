# Build
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o main .

# Prod
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main /main
COPY elements.json /app
RUN apk add --no-cache libc6-compat
EXPOSE 8080
ENTRYPOINT ["/main"]
