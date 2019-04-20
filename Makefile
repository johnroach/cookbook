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
	docker run -p 8080:8080 --mount type=bind,source=$(shell PWD)/config,target=/config cookbook:$(GIT_COMMIT)

docker-tag-gcloud: docker-build
	@printf '\n\n\n\e[1;34m%-6s\e[m\n\n\n' "RETAGGING CONTAINER FOR GCLOUD"
	docker tag cookbook:$(GIT_COMMIT) us.gcr.io/red-fa1ebe00/cookbook:$(GIT_COMMIT)

docker-push-gcloud: docker-tag-gcloud
	@printf '\n\n\n\e[1;34m%-6s\e[m\n\n\n' "PUSHING CONTAINER TO GCLOUD"
	docker push us.gcr.io/red-fa1ebe00/cookbook:$(GIT_COMMIT)

deploy-local-k8s-all: get-version docker-build
	@printf '\n\n\n\e[1;34m%-6s\e[m\n\n\n' "DEPLOYING TO LOCAL K8S"
	kubectx docker-desktop
	rm -rf deployment/deployment.yml.tmp1
	sed 's/@@IMAGE_NAME@@/cookbook/g' deployment/deployment.yml > deployment/deployment.yml.tmp
	sed "s/@@VERSION@@/$(GIT_COMMIT)/g" deployment/deployment.yml.tmp > deployment/deployment.yml.tmp1
	rm -rf deployment/deployment.yml.tmp
	kubectl create configmap cookbook-config --from-file=config/dev.yaml || kubectl create configmap cookbook-config --from-file config/dev.yaml -o yaml --dry-run | kubectl replace -f -
	kubectl apply -f deployment/deployment.yml.tmp1

deploy-local-k8s: get-version
	@printf '\n\n\n\e[1;34m%-6s\e[m\n\n\n' "DEPLOYING TO LOCAL K8S"
	kubectx docker-desktop
	rm -rf deployment/deployment.yml.tmp1
	sed 's/@@IMAGE_NAME@@/cookbook/g' deployment/deployment.yml > deployment/deployment.yml.tmp
	sed "s/@@VERSION@@/$(GIT_COMMIT)/g" deployment/deployment.yml.tmp > deployment/deployment.yml.tmp1
	rm -rf deployment/deployment.yml.tmp
	kubectl create configmap cookbook-config --from-file=config/dev.yaml || kubectl create configmap cookbook-config --from-file config/dev.yaml -o yaml --dry-run | kubectl replace -f -
	kubectl apply -f deployment/deployment.yml.tmp1
