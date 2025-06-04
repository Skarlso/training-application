IMAGE_REPOSITORY = quay.io/kubermatic-labs
APPLICATION_NAME = training-application
BUILD_VERSION = 4.0.0
BUILD_VERSION_A = ${BUILD_VERSION}-A
BUILD_VERSION_B = ${BUILD_VERSION}-B
BUILD_VERSION_DISTROLESS = ${BUILD_VERSION}-distroless

.PHONY: update-dependencies
update-dependencies: 
	cd src && go get -u
	cd src && go mod tidy

.PHONY: lint
lint:
	cd src && golangci-lint run --timeout=5m -v 

.PHONY: build
build: 
	cd src && go build -o ../${APPLICATION_NAME} .

.PHONY: run
run: build
	./${APPLICATION_NAME}

.PHONY: docker-lint
docker-lint: 
	hadolint Dockerfile
	hadolint Dockerfile-A
	hadolint Dockerfile-B
	hadolint Dockerfile-distroless

.PHONY: docker-build
docker-build: build
	docker build -t ${IMAGE_REPOSITORY}/${APPLICATION_NAME}:${BUILD_VERSION} .

.PHONY: docker-run
docker-run: docker-build
	docker run -it --rm -p 8080:8080 -m=10m --cpus=".5" --name ${APPLICATION_NAME} ${IMAGE_REPOSITORY}/${APPLICATION_NAME}:${BUILD_VERSION}

.PHONY: docker-build-all
docker-build-all: 
	docker build -t ${IMAGE_REPOSITORY}/${APPLICATION_NAME}:${BUILD_VERSION} .
	docker build -f Dockerfile-A -t ${IMAGE_REPOSITORY}/${APPLICATION_NAME}:${BUILD_VERSION_A} .
	docker build -f Dockerfile-B -t ${IMAGE_REPOSITORY}/${APPLICATION_NAME}:${BUILD_VERSION_B} .
	docker build -f Dockerfile-distroless -t ${IMAGE_REPOSITORY}/${APPLICATION_NAME}:${BUILD_VERSION_DISTROLESS} .

.PHONY: docker-push
docker-push: 
	docker buildx build --push --platform linux/arm64,linux/amd64 --tag ${IMAGE_REPOSITORY}/${APPLICATION_NAME}:${BUILD_VERSION} .

.PHONY: docker-push-all
docker-push-all: 
	docker buildx build --push --platform linux/arm64,linux/amd64 --tag ${IMAGE_REPOSITORY}/${APPLICATION_NAME}:${BUILD_VERSION} .
	docker buildx build --push --platform linux/arm64,linux/amd64 -f Dockerfile-A --tag ${IMAGE_REPOSITORY}/${APPLICATION_NAME}:${BUILD_VERSION_A} .
	docker buildx build --push --platform linux/arm64,linux/amd64 -f Dockerfile-B --tag ${IMAGE_REPOSITORY}/${APPLICATION_NAME}:${BUILD_VERSION_B} .
	docker buildx build --push --platform linux/arm64,linux/amd64 -f Dockerfile-distroless --tag ${IMAGE_REPOSITORY}/${APPLICATION_NAME}:${BUILD_VERSION_DISTROLESS} .
