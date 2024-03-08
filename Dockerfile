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

# Copy gateway binary
COPY --from=builder /gateway/bin/gateway-amd64-linux /bin/gateway-service

# Copy migrations directory
COPY --from=builder /gateway/internal/storage/migrations /migrations

RUN ls /bin/gateway-service
WORKDIR /bin
ENTRYPOINT [ "./gateway-service" ]