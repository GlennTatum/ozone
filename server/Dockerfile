FROM golang:1.24 AS builder
WORKDIR /app
COPY . .
RUN go build -o server main.go

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8000
ENV APP_ENV=production
ENV CASSANDRA_URL=localhost:9042
ENV REDIS_URL=localhost:6379
ENTRYPOINT ["./server"]
