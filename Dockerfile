# Builder image go
FROM golang:1.22.0 AS builder

ARG appVersion
ARG commitHash

ENV VERSION=$appVersion
ENV COMMIT_HASH=$commitHash

# Build gateway binary with Go
ENV GOPATH /opt/go

RUN mkdir -p /gateway
WORKDIR /gateway
COPY . /gateway
RUN go mod tidy && make build-linux

# Runnable image
FROM alpine:3.19
ARG appVersion
ARG commitHash
ENV VERSION=$appVersion
ENV COMMIT_HASH=$commitHash

RUN apk add --no-cache bash

# Copy gateway binary
COPY --from=builder /gateway/bin/gateway-amd64-linux /bin/gateway-service

# Copy migrations directory
COPY --from=builder /gateway/internal/storage/migrations /migrations

# Add wait-for-it
COPY --from=builder /gateway/wait-for-it.sh /bin/wait-for-it.sh
RUN chmod +x /bin/wait-for-it.sh

RUN ls /bin/gateway-service
WORKDIR /bin
ENTRYPOINT ["./wait-for-it.sh", "database:3306", "--", "./gateway-service"]