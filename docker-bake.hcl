group "default" {
  targets = ["production"]
}

target "production" {
  context = "."
  dockerfile = "docker/prod.Dockerfile"
  args = {
    ALPINE_VERSION = "latest"
    GO_VERSION = "1.24.1"
  }
  tags = ["zhangyumao:latest"]
}
