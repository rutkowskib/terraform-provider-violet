terraform {
  required_providers {
    violet = {
      source = "hashicorp.com/edu/violet"
    }
  }
}

provider "violet" {
  username   = var.username
  password   = var.password
  app_id     = var.app_id
  app_secret = var.app_secret
}

resource "violet_webhook" "example" {
  event           = "OFFER_UPDATED"
  remote_endpoint = "https://test.com/"
}
