
#Build binary
FROM golang:1.16-alpine as builder
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

WORKDIR /workspace

# Copy sources
COPY ./ .

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o csv-exporter ./cmd/csv-service/main.go

# Executuion Image
FROM scratch
WORKDIR /

# Import from builder.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd

COPY --from=builder /workspace/csv-exporter .
ENTRYPOINT ["/csv-exporter"]
