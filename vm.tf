###############################################################################
# Server
###############################################################################

resource "google_compute_instance" "game_server" {
  project      = google_project.project.name
  name         = var.server_name
  machine_type = var.machine_type
  zone         = var.default_zone

  allow_stopping_for_update = true

  boot_disk {
    device_name = "boot"
    auto_delete = true
    initialize_params {
      image = data.google_compute_image.ubuntu_image.self_link
    }
  }

  attached_disk {
    device_name = "game-storage"
    source      = google_compute_disk.game_storage.self_link
  }

  network_interface {
    network    = google_compute_network.network.self_link
    subnetwork = google_compute_subnetwork.subnet.self_link
    access_config {
      nat_ip = google_compute_address.static_external_ip.address
    }
  }

  tags = local.game_server_network_tags

  metadata = {
    startup-script  = local.game_server_startup_script
    shutdown-script = local.game_server_shutdown_script
  }

  service_account {
    scopes = [
      "cloud-platform",
      "logging-write"
    ]
  }
}

data "google_compute_image" "ubuntu_image" {
  family  = "ubuntu-2004-lts"
  project = "ubuntu-os-cloud"
}

###############################################################################
# Disk
###############################################################################

resource "google_compute_disk" "game_storage" {
  project = google_project.project.name
  name    = "game-storage"
  type    = "pd-ssd"
  zone    = var.default_zone
  size    = var.disk_size

  lifecycle {
    prevent_destroy = true
  }
}

resource "google_compute_disk_resource_policy_attachment" "backup_policy_attachment" {
  project = google_project.project.name
  name    = google_compute_resource_policy.backup_policy.name
  disk    = google_compute_disk.game_storage.name
  zone    = var.default_zone
}

resource "google_compute_resource_policy" "backup_policy" {
  project = google_project.project.name
  name    = "backup-policy"
  region  = var.default_region

  snapshot_schedule_policy {
    schedule {
      daily_schedule {
        days_in_cycle = 1
        start_time    = "05:00"
      }
    }
    retention_policy {
      max_retention_days = 10
    }
    snapshot_properties {
      storage_locations = [var.default_region]
      guest_flush       = true
    }
  }
}
