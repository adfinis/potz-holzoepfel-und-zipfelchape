project = "caasperli"

app "caasperli" {
  labels = {
    "service" = "caasperli",
    "env" = "dev"
  }

  build {
    use "pack" {}
  }

  deploy {
    use "docker" {
        service_port = 8080
    }
  }
}

