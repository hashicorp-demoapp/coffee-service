project = "hashicups"

app "coffee-service" {
  labels = {
    "service" = "coffee-service",
    "env" = "dev"
  }

  build {
    use "docker" {}
    registry {
      use "docker" {
        image = "gcr.io/consul-k8s-<redacted>/coffee-service"
        tag = "v3.0.0"
      }
    }
  }

  deploy {
    use "kubernetes" {
      probe_path = "/health"
      service_port = 9090
      annotations = {
        "prometheus.io/scrape" = "true",
        "prometheus.io/port" = "9102",
        "consul.hashicorp.com/connect-inject" = "true",
        "consul.hashicorp.com/service-meta-version" = "v3",
        "consul.hashicorp.com/service-tags" = "api"
      }
      static_environment = {
        LOG_FORMAT = "text",
        LOG_LEVEL = "INFO",
        BIND_ADDRESS = "localhost:9090",
        METRICS_ADDRESS = "localhost:9102",
        VERSION = "v3"
      }
    }
  }

  release {
    use "kubernetes" {
    }
  }
}