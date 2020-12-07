# coffee-service

The coffee service illustrates how to extract a microservice out of a monolith at both a code level, and by leveraging the service mesh features of Consul. It also illustrates how to instrument code for distributed tracing with Jaeger.

## Code changes

All coffee related functionality from the product-api has been extracted over to this microservice. The service has been instrumented with the go-hckit tracing middleware. There are multiple versions of the service with different featuresets by version.

- v1 provides all the core db lookup logic and a simple Jaeger instrumentation implementation
- v2 improves the instrumentation position by leveraging [jmoiron/sqlx](https://github.com/jmoiron/sqlx) so that users can optionally deploy SQL_TRACE_ENABLED builds for query tracing via Jaeger
- v3 improves the implementation by converting the service to a search node using local data, and thus sidestepping the database calls entirely

## Included configuration

- coffee-service-v1.yaml - Deployment for v1 of the service
- coffee-service-v2.yaml - Deployment for v2 of the service
- coffee-service-v3.yaml - Deployment for v3 of the service
- service-router.yaml - Consul config entry that can be applied to re-route traffic away from products-api based on route
- service-splitter.yaml - Consul config entry that can be applied define canary deployment percentages
- service-resolver.yaml - Consul config entry that can be applied to resolve canary deployment splits

## Running locally

First build the local docker image.

`make build_docker_dev`

Then run the image. You can supply whatever port, version, log level, and name you like

`docker run -d -p 9090:9090 --env BIND_ADDRESS=localhost:9090 --env VERSION=v3 --env LOG_LEVEL=DEBUG --name=coffee-service hashicorpdemoapp/coffee-service:devlocal`
