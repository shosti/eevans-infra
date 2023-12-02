terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.29.0"
    }

    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 3.4.0"
    }
  }
}
