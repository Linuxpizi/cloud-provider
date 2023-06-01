.PHONY: tidy
tidy:
	@go mod tidy

.PHONY: build-image
build-image: tidy
	docker buildx build -t linuxpizi/k8s-client-in-cluster:latest -f Dockerfile .