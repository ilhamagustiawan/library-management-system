FROM golang:1.23-alpine AS builder
WORKDIR /src
RUN apk add --no-cache ca-certificates
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/auth-service .

FROM alpine:3.21
RUN apk add --no-cache ca-certificates && addgroup -S app && adduser -S -G app app
WORKDIR /app
COPY --from=builder /out/auth-service /app/auth-service
COPY --from=builder /src/db/migrations /app/db/migrations
USER app
EXPOSE 8081
CMD ["/app/auth-service", "serve"]
