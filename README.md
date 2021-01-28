# coffee-service
Coffee API written in Go

[![CircleCI](https://circleci.com/gh/hashicorp-demoapp/coffee-service.svg?style=svg)](https://circleci.com/gh/hashicorp-demoapp/coffee-service)

Docker Image: [https://hub.docker.com/r/hashicorpdemoapp/coffee-service](https://hub.docker.com/r/hashicorpdemoapp/coffee-service)

The coffee service illustrates how to extract a microservice out of a monolith.

## Code changes

All coffee related functionality from the product-api has been extracted over to this microservice. The service has been
instrumented with the go-hckit tracing middleware. There are multiple versions of the service with different featuresets
by version.

- v1 provides all the core db lookup logic and a simple Jaeger instrumentation implementation
- v2 improves the instrumentation position by leveraging [jmoiron/sqlx](https://github.com/jmoiron/sqlx) so that users
  can optionally deploy `SQL_TRACE_ENABLED` builds for query tracing via Jaeger
- v3 improves the implementation by converting the service to a search node using in memory data, and thus sidestepping
  the database calls entirely

## Included Kubernetes configuration

- coffee-service-v1.yaml - Deployment for v1 of the service
- coffee-service-v2.yaml - Deployment for v2 of the service
- coffee-service-v3.yaml - Deployment for v3 of the service

## Included Waypoint configuration

Waypoint support is under active development. The `waypoint.hcl` file in the project root is not yet stable. Check back
for updates.

## Running locally

First build the local docker image.

`make build_docker_dev`

Then run the image. You can supply whatever port, version, log level, and name you like

`docker run -d -p 9090:9090 --env BIND_ADDRESS=localhost:9090 --env VERSION=v3 --env LOG_LEVEL=DEBUG --name=coffee-service hashicorpdemoapp/coffee-service:devlocal`
