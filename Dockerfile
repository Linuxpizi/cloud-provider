# syntax=docker/dockerfile:1
FROM golang:1.20 AS builder
WORKDIR /build
COPY . .
RUN go mod tidy && \
    CGO_ENABLED=0 go build -o ucloudstack-cloud-controller-manager ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /bin
COPY --from=builder /build ./
CMD ["/bin/ucloudstack-cloud-controller-manager"]
