FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /out/prosel-api ./cmd/api

FROM alpine:3.22

WORKDIR /app
COPY --from=builder /out/prosel-api /usr/local/bin/prosel-api
COPY migrations ./migrations
EXPOSE 8080
CMD ["prosel-api"]
