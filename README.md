# cookbook

## Development

Download and install [protobuf](https://github.com/protocolbuffers/protobuf/releases)

Run the following commands

```bash
go get -u google.golang.org/grpc
go get -u github.com/golang/protobuf/protoc-gen-go
```

Generating code via proto files

```bash
mkdir -p proto-gen/helloworld; protoc -I model/proto model/proto/helloworld.proto --go_out=plugins=grpc:proto-gen/helloworld
mkdir -p proto-gen/health; protoc -I model/proto model/proto/health.proto --go_out=plugins=grpc:proto-gen/health
```

## References

### Technical References
- [Health Probe for GRPC](https://github.com/grpc-ecosystem/grpc-health-probe)
- [Testing GRPC](https://github.com/grpc/grpc-go/blob/master/Documentation/gomock-example.md)

### Birds of the Feather
- [How we build GRPC services at Namely?](https://medium.com/namely-labs/how-we-build-grpc-services-at-namely-52a3ae9e7c35)
- [Awesome GRPC - GRPC Ecosystem](https://github.com/grpc-ecosystem/awesome-grpc)
