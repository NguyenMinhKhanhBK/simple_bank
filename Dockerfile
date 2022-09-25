# Build stage
FROM golang:1.19.1-alpine3.16 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add curl
ENV GOLANG_MIGRATE_VERSION=v4.15.2
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/${GOLANG_MIGRATE_VERSION}/migrate.linux-amd64.tar.gz | tar xvz

# Run stage
FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate .
COPY app.env .
COPY db/migration ./migration/
COPY start.sh .
COPY wait-for.sh .

EXPOSE 8080
CMD ["/app/main"]
ENTRYPOINT [ "/app/start.sh" ]
