# Build stage
FROM  golang:1.18-alpine3.16 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/ .
COPY . .
COPY .env .
COPY start.sh .
COPY wait-for.sh .


CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]
