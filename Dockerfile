FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o testum app/cmd/testum/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/testum .

COPY --from=builder /app/app/migrations ./migrations

EXPOSE 8080

CMD ["./testum"]
