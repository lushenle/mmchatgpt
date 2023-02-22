# Build the tlssecretmanager binary
FROM golang:1.19 as builder

WORKDIR /workspace

# Install upx for compress binary file
RUN apt update && apt install -y upx

# Copy the go source
COPY . .

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GO111MODULE=on

# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Build and compression
RUN go build -a -installsuffix cgo -ldflags="-s -w" -o mmchatgpt main.go \
    && upx mmchatgpt

# build tls
FROM alpine-glibc:glibc-2.34 as final
WORKDIR /
COPY --from=builder /workspace/tlssecretmanager .
#COPY config.yaml .

ENTRYPOINT ["/mmchatgpt"]
