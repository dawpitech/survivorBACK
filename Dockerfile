# Base image for Go
FROM golang:1.23.8-alpine AS base
WORKDIR /app

# Dependency stage: only download modules if go.mod/go.sum change
FROM base AS deps
COPY go.mod go.sum ./
RUN go mod download

# Builder stage: copy sources and build
FROM base AS builder
COPY --from=deps /go/pkg /go/pkg
COPY . .
RUN go build -o api .

# Final minimal image
FROM alpine:latest AS release
WORKDIR /app
COPY --from=builder /app/api .
EXPOSE 24680
CMD ["./api"]