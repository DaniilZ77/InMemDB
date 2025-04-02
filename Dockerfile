FROM golang:1.24.2 AS builder
WORKDIR /app
COPY . .
ENV CGO_ENABLED=0
RUN go build -o ./bin/db ./cmd/db

FROM alpine:3.21.3
WORKDIR /app
COPY --from=builder ./app/bin ./bin

ENTRYPOINT [ "./bin/db" ]