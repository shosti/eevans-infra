terraform {
  required_providers {
    minio = {
      source  = "aminueza/minio"
      version = "~> 1.20.0"
    }

    healthchecksio = {
      source  = "kristofferahl/healthchecksio"
      version = "~> 1.10.0"
    }
  }
}
