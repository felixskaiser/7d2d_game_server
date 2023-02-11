terraform {
  backend "gcs" {
    bucket  = "terraform_state_private"
    prefix  = "state/game-server-7d2d-felix"
  }
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "4.52.0"
    }
  }
}

provider "google" {
}