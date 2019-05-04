
deploy-local-k8s-all: get-version docker-build
	@printf '\n\n\n\e[1;34m%-6s\e[m\n\n\n' "DEPLOYING TO LOCAL K8S"
	kubectx docker-desktop
	rm -rf deployment/deployment.yml.tmp1
	sed 's/@@IMAGE_NAME@@/cookbook/g' deployment/local_deployment.yml > deployment/deployment.yml.tmp
	sed "s/@@VERSION@@/$(GIT_COMMIT)/g" deployment/deployment.yml.tmp > deployment/deployment.yml.tmp1
	rm -rf deployment/deployment.yml.tmp
	kubectl create configmap cookbook-config --from-file=config/dev.yaml || kubectl create configmap cookbook-config --from-file config/dev.yaml -o yaml --dry-run | kubectl replace -f -
	kubectl apply -f deployment/deployment.yml.tmp1

deploy-local-k8s: get-version
	@printf '\n\n\n\e[1;34m%-6s\e[m\n\n\n' "DEPLOYING TO LOCAL K8S"
	kubectx docker-desktop
	rm -rf deployment/deployment.yml.tmp1
	sed 's/@@IMAGE_NAME@@/cookbook/g' deployment/local_deployment.yml > deployment/deployment.yml.tmp
	sed "s/@@VERSION@@/$(GIT_COMMIT)/g" deployment/deployment.yml.tmp > deployment/deployment.yml.tmp1
	rm -rf deployment/deployment.yml.tmp
	kubectl create configmap cookbook-config --from-file=config/dev.yaml || kubectl create configmap cookbook-config --from-file config/dev.yaml -o yaml --dry-run | kubectl replace -f -
	kubectl apply -f deployment/deployment.yml.tmp1
