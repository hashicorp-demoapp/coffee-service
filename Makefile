CONTAINER_NAME=hashicorpdemoapp/coffee-service
DB_CONTAINER_NAME=hashicorpdemoapp/product-api-db
CONTAINER_VERSION=v0.0.1

test_functional:
	shipyard run ./blueprint
	cd ./functional_tests && go test -v -run.test true ./..
	shipyard destroy

build_linux:
	CGO_ENABLED=0 GOOS=linux go build -o ./bin/coffee-service

build_docker: build_linux
	docker build -t ${CONTAINER_NAME}:${CONTAINER_VERSION} .

build_docker_dev: build_linux
	docker build -t ${CONTAINER_NAME}:devlocal .

# docker run -d -p 9090:9090 --env BIND_ADDRESS=localhost:9090 --env VERSION=v3 --env LOG_LEVEL=DEBUG --name=coffee-service hashicorpdemoapp/coffee-service:devlocal

push_docker: build_docker
	docker push ${CONTAINER_NAME}:${CONTAINER_VERSION}