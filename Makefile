get-version:
	$(eval GIT_COMMIT = $(shell git rev-parse --short=7 HEAD))

dep:
	@printf '\n\n\n\e[1;34m%-6s\e[m\n\n\n' "DOWNLOADING DEPENDENCIES"
	go mod download

test: dep
	@printf '\n\n\n\e[1;34m%-6s\e[m\n\n\n' "RUNNING TESTS"
	go test -v ./...

build-dev: test
	@printf '\n\n\n\e[1;34m%-6s\e[m\n\n\n' "BUILDING"
	go build

run-dev: test
	@printf '\n\n\n\e[1;34m%-6s\e[m\n\n\n' "RUNNING APPLICATION"
	go run main.go -e dev

docker-build: get-version test
	@printf '\n\n\n\e[1;34m%-6s\e[m\n\n\n' "BUILDING DOCKER CONTAINER"
	docker build -t cookbook:$(GIT_COMMIT) --build-arg VERSION=$(GIT_COMMIT) .

docker-dev-run: get-version
	@printf '\n\n\n\e[1;34m%-6s\e[m\n\n\n' "RUNNING DOCKER CONTAINER"
	docker run -p 8080:8080 --mount type=bind,source=$(shell PWD)/config,target=/config cookbook:$(GIT_COMMIT)

docker-dev-run-all: docker-build
	@printf '\n\n\n\e[1;34m%-6s\e[m\n\n\n' "RUNNING DOCKER CONTAINER"
	docker run -p 8080:8080 --mount type=bind,source=$(shell PWD)/config,target=/config cookbook:$(GIT_COMMIT