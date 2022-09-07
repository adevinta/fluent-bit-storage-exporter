FROM golang:1.17 as builder
WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

# Copy the go source
COPY pkg/ pkg/
COPY vendor/ vendor/
COPY cmd/ cmd/
# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -mod=vendor -a -o fluent-bit-storage-exporter cmd/main.go

FROM gcr.io/distroless/static:nonroot
MAINTAINER CPR Team <common-platform-runtime@adevinta.com>
WORKDIR /
COPY --from=builder /workspace/fluent-bit-storage-exporter .
USER nonroot:nonroot
ENTRYPOINT ["/fluent-bit-storage-exporter"]
