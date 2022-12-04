terraform {
  required_providers {
    commercelayer = {
      version = ">= 0.0.1"
      source  = "incentro-dc/commercelayer"
    }
  }
}

provider "commercelayer" {
  client_id     = "<client-id>"
  client_secret = "<client-secret>"
  api_endpoint  = "<api-endpoint>"
  auth_endpoint = "<auth-endpoint>"
}