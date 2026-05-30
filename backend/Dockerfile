FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /server ./cmd/server
RUN CGO_ENABLED=0 go build -o /seed ./cmd/seed

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
COPY --from=builder /server /server
COPY --from=builder /seed /seed
COPY db/migrations /db/migrations
EXPOSE 8080
ENTRYPOINT ["/server"]
