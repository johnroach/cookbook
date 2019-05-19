FROM golang:1.12 as builder

# Pulling in the health probe
# You can read more about it here https://github.com/grpc-ecosystem/grpc-health-probe
RUN GRPC_HEALTH_PROBE_VERSION=v0.2.2 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

# I can't believe we still need this stupid env var for 1.12
ENV GO111MODULE=on

# Setting up cookbook build
COPY . $GOPATH/src/github.com/johnroach/cookbook/
WORKDIR $GOPATH/src/github.com/johnroach/cookbook/
RUN go mod download

# Running cookbook build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/cookbook

# Starting on Scratch
FROM scratch

ARG VERSION

# Moving needed binaries to
COPY --from=builder /go/bin/cookbook /go/bin/cookbook
COPY --from=builder /bin/grpc_health_probe /go/bin/grpc_health_probe
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENV COOKBOOK_VERSION=${VERSION}

ENTRYPOINT ["/go/bin/cookbook"]