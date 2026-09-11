FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download && go build -o server ./cmd/server

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 50051
CMD ["./server"]