# cookbook

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/45fb61c8af604b7cb45679ad240fa4e0)](https://app.codacy.com/app/johnroach1985/cookbook?utm_source=github.com&utm_medium=referral&utm_content=johnroach/cookbook&utm_campaign=Badge_Grade_Dashboard) [![CodeFactor](https://www.codefactor.io/repository/github/johnroach/cookbook/badge)](https://www.codefactor.io/repository/github/johnroach/cookbook) [![codecov](https://codecov.io/gh/johnroach/cookbook/branch/master/graph/badge.svg)](https://codecov.io/gh/johnroach/cookbook) [![Build Status](https://travis-ci.org/johnroach/cookbook.svg?branch=master)](https://travis-ci.org/johnroach/cookbook) [![Go Report Card](https://goreportcard.com/badge/github.com/johnroach/cookbook)](https://goreportcard.com/report/github.com/johnroach/cookbook)

## Development

Install magefile:

```$bash
go get -u -d github.com/magefile/mage
cd $GOPATH/src/github.com/magefile/mage
go run bootstrap.go
```

or do:

```$bash
git clone https://github.com/magefile/mage
cd mage
go run bootstrap.go
```

Download and install [protobuf](https://github.com/protocolbuffers/protobuf/releases)

Run the following commands



Generating code via proto files ( TODO: Move this work to magefile)

```bash
mkdir -p proto-gen/helloworld; protoc -I model/proto model/proto/helloworld.proto --go_out=plugins=grpc:proto-gen/helloworld
mkdir -p proto-gen/health; protoc -I model/proto model/proto/health.proto --go_out=plugins=grpc:proto-gen/health
```

## References

### Technical References

-[Health Probe for GRPC](https://github.com/grpc-ecosystem/grpc-health-probe)
-[Testing GRPC](https://github.com/grpc/grpc-go/blob/master/Documentation/gomock-example.md)
-[Project Layout](https://github.com/golang-standards/project-layout)

### Birds of the Feather

-[How we build GRPC services at Namely?](https://medium.com/namely-labs/how-we-build-grpc-services-at-namely-52a3ae9e7c35)
-[Awesome GRPC - GRPC Ecosystem](https://github.com/grpc-ecosystem/awesome-grpc)
