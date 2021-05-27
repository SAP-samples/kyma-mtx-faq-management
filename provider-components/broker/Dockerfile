
#Build binary
FROM golang:1.14-alpine as builder
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

WORKDIR /workspace
# Copy the go source
COPY main.go main.go
COPY pkg/ pkg/

# Copy the Go Modules manifests
COPY go.mod go.mod

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o broker main.go

# Executuion Image
FROM scratch
WORKDIR /

# Import from builder.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd

COPY --from=builder /workspace/broker .
ENTRYPOINT ["/broker"]
