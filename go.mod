module github.com/hashicorp-demoapp/coffee-service

go 1.14

require (
	contrib.go.opencensus.io/integrations/ocsql v0.1.6
	github.com/cucumber/godog v0.10.0
	github.com/cucumber/messages-go/v10 v10.0.3
	github.com/fsnotify/fsnotify v1.4.9
	github.com/gorilla/mux v1.8.0
	github.com/hashicorp-demoapp/go-hckit v0.0.1
	github.com/hashicorp-demoapp/product-api-go v0.0.12
	github.com/hashicorp/consul v1.8.4
	github.com/hashicorp/go-hclog v0.14.1
	github.com/jmoiron/sqlx v1.2.0
	github.com/lib/pq v1.2.0
	github.com/nicholasjackson/env v0.6.0
	github.com/opentracing/opentracing-go v1.2.0
	github.com/stretchr/testify v1.6.1
	github.com/uber/jaeger-client-go v2.25.0+incompatible // indirect
	github.com/uber/jaeger-lib v2.2.0+incompatible // indirect
)

// replace github.com/DerekStrickland/learn-consul-jaeger/go-hckit => /Users/derekstrickland/code/DerekStrickland/learn-consul-jaeger/go-hckit
