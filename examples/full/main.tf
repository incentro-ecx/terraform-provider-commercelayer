terraform {
  required_providers {
    commercelayer = {
      version = ">= 0.0.1"
      source  = "incentro/commercelayer"
    }
  }
}

provider "commercelayer" {
  client_id     = var.client_id
  client_secret = var.client_secret
  api_endpoint  = var.api_endpoint
  auth_endpoint = var.auth_endpoint
}