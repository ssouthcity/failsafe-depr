terraform {
  cloud {
    organization = "southcity"

    workspaces {
      name = "failsafe"
    }
  }

  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
  }
}

variable "do_token" {
  type = string
}

provider "digitalocean" {
  token = var.do_token
}

data "digitalocean_ssh_key" "neptr" {
  name = "NEPTR"
}

resource "digitalocean_droplet" "failsafe" {
  image    = "ubuntu-20-04-x64"
  name     = "failsafe-1"
  region   = "fra1"
  size     = "s-1vcpu-1gb"
  ssh_keys = [data.digitalocean_ssh_key.neptr.id]
  user_data = templatefile("${path.module}/user-data.sh", {
    config_file_content  = file("${path.module}/config.json")
    service_file_content = file("${path.module}/failsafe.service")
  })
}

output "droplet_id" {
  value = digitalocean_droplet.failsafe.id
}
