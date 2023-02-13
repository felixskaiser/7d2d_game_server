###############################################################################
# Generate password required to join game server as player
###############################################################################

resource "random_password" "server_password" {
  length  = 8
  special = false
}

resource "google_secret_manager_secret" "server_password_secret" {
  project   = google_project.project.name
  secret_id = "server-password"

  replication {
    user_managed {
      replicas {
        location = var.default_region
      }
    }
  }
  depends_on = [google_project_service.secret_manager_service]
}

resource "google_secret_manager_secret_version" "server_password_secret_version" {
  secret      = google_secret_manager_secret.server_password_secret.id
  secret_data = random_password.server_password.result
  depends_on  = [google_project_service.secret_manager_service]
}

###############################################################################
# Generate password for game server control panel
###############################################################################

resource "random_password" "control_panel_password" {
  length  = 12
  special = false
}

resource "google_secret_manager_secret" "control_panel_password_secret" {
  project   = google_project.project.name
  secret_id = "control-panel-password"

  replication {
    user_managed {
      replicas {
        location = var.default_region
      }
    }
  }
  depends_on = [google_project_service.secret_manager_service]
}

resource "google_secret_manager_secret_version" "control_panel_password_secret_version" {
  secret      = google_secret_manager_secret.control_panel_password_secret.id
  secret_data = random_password.control_panel_password.result
  depends_on  = [google_project_service.secret_manager_service]
}

###############################################################################
# Generate password for telnet server
###############################################################################

resource "random_password" "telnet_password" {
  length  = 12
  special = false
}

resource "google_secret_manager_secret" "telnet_password_secret" {
  project   = google_project.project.name
  secret_id = "telnet-password"

  replication {
    user_managed {
      replicas {
        location = var.default_region
      }
    }
  }
  depends_on = [google_project_service.secret_manager_service]
}

resource "google_secret_manager_secret_version" "telnet_secret_version" {
  secret      = google_secret_manager_secret.telnet_password_secret.id
  secret_data = random_password.telnet_password.result
  depends_on  = [google_project_service.secret_manager_service]
}

###############################################################################
# Generate password for server manager functions
###############################################################################

resource "random_password" "server_manager_password" {
  length  = 12
  special = false
}

resource "google_secret_manager_secret" "server_manager_password_secret" {
  project   = google_project.project.name
  secret_id = "server-manager-password"

  replication {
    user_managed {
      replicas {
        location = var.default_region
      }
    }
  }
  depends_on = [google_project_service.secret_manager_service]
}

resource "google_secret_manager_secret_version" "server_manager_secret_version" {
  secret      = google_secret_manager_secret.server_manager_password_secret.id
  secret_data = random_password.server_manager_password.result
  depends_on  = [google_project_service.secret_manager_service]
}
