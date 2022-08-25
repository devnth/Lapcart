# Build stage
FROM  golang:1.18-alpine3.16 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/main .
COPY . .
COPY .env .
COPY start.sh .
COPY wait-for.sh .


CMD [ "/app/main.go" ]
ENTRYPOINT [ "/app/start.sh" ]
