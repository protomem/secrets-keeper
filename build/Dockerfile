FROM golang:alpine AS builder

RUN apk --no-cache add gcc libc-dev sqlite

WORKDIR /app

COPY go.* .
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -v -o ./build/ ./cmd/secrets-keeper



FROM alpine:latest

RUN apk --no-cache add ca-certificates
RUN apk --no-cache add sqlite

WORKDIR /app

COPY --from=builder /app/configs/ ./configs/
COPY --from=builder /app/build/ ./build/
RUN mkdir -p /app/data

CMD [ "./build/secrets-keeper" ]
