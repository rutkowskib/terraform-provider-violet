terraform {
  required_providers {
    violet = {
      source = "rutkowskib/violet"
    }
  }
}

provider "violet" {
  username   = var.username
  password   = var.password
  app_id     = var.app_id
  app_secret = var.app_secret
}