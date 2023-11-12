terraform {
  required_providers {
    minio = {
      source  = "aminueza/minio"
      version = "~> 1.18.0"
    }

    healthchecksio = {
      source  = "kristofferahl/healthchecksio"
      version = "~> 1.9.0"
    }
  }
}
