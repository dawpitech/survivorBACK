FROM golang:1.23.8-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o api .

FROM alpine:latest AS release
WORKDIR /app
COPY --from=builder /app/api .
EXPOSE 24680
CMD [ "./api" ]