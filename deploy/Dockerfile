FROM golang:1.23.6-alpine AS builder
WORKDIR /build
COPY . .
RUN go build -o main ./src/cmd/main.go

FROM alpine:3.21 AS runner
WORKDIR /app
COPY --from=builder /build/main .
CMD ["./main"]
