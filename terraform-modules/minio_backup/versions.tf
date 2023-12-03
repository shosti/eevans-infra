terraform {
  required_providers {
    minio = {
      source  = "aminueza/minio"
      version = "~> 2.0.0"
    }

    healthchecksio = {
      source  = "kristofferahl/healthchecksio"
      version = "~> 1.9.0"
    }
  }
}
