FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /main ./cmd/app
RUN CGO_ENABLED=0 GOOS=linux go build -o /migrate ./cmd/migrate

FROM alpine:latest AS runner

WORKDIR /

COPY --from=builder /main /main
COPY --from=builder /migrate /migrate
COPY --from=builder /app/migrations /migrations
COPY --from=builder /app/seeds /seeds
COPY --from=builder /app/web /web

EXPOSE 8080
CMD ["/main"]

