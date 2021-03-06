CONTAINER_NAME=hashicorpdemoapp/coffee-service
DB_CONTAINER_NAME=hashicorpdemoapp/coffee-service
CONTAINER_VERSION=v0.0.2

test_functional:
	cd ./functional_tests && go test -v -run.test true ./..

build_linux:
	CGO_ENABLED=0 GOOS=linux go build -o ./bin/coffee-service

build_docker: build_linux
	docker build -t ${CONTAINER_NAME}:${CONTAINER_VERSION} .

build_docker_dev: build_linux
	docker build -t ${CONTAINER_NAME}:devlocal .

push_docker: build_docker
	docker push ${CONTAINER_NAME}:${CONTAINER_VERSION}
