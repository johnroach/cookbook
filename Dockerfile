FROM golang:1.12-stretch as builder
# i can't believe we still need this stupid env var...
ENV GO111MODULE=on
COPY . $GOPATH/src/github.com/johnroach/cookbook/
WORKDIR $GOPATH/src/github.com/johnroach/cookbook/
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/cookbook


FROM scratch

ARG VERSION

COPY --from=builder /go/bin/cookbook /go/bin/cookbook
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENV COOKBOOK_VERSION=${VERSION}

ENTRYPOINT ["/go/bin/cookbook"]