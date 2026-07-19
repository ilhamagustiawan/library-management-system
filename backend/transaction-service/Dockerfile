FROM golang:1.23-alpine AS builder
WORKDIR /src
RUN apk add --no-cache ca-certificates
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/transaction-service .

FROM alpine:3.21
RUN apk add --no-cache ca-certificates && addgroup -S app && adduser -S -G app app
WORKDIR /app
COPY --from=builder /out/transaction-service /app/transaction-service
COPY --from=builder /src/db /app/db
COPY --from=builder /src/docs /app/docs
USER app
EXPOSE 8080
CMD ["/app/transaction-service", "serve"]
