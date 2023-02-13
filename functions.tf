resource "google_storage_bucket" "cloud_function_staging_bucket" {
  name                        = "${google_project.project.name}-gcf-source" # Every bucket name must be globally unique
  project                     = google_project.project.name
  location                    = "EU"
  uniform_bucket_level_access = true
  public_access_prevention    = "enforced"

  autoclass {
    enabled = true
  }

  versioning {
    enabled = false
  }
}

###############################################################################
# Cloud Function 'status'
###############################################################################

resource "google_cloudfunctions2_function" "function_status" {
  name     = "status"
  project  = google_project.project.name
  location = var.default_region

  build_config {
    runtime     = "go119"
    entry_point = "Entrypoint"

    source {
      storage_source {
        bucket = google_storage_bucket.cloud_function_staging_bucket.name
        object = google_storage_bucket_object.function_status_source.name
      }
    }
  }

  service_config {
    available_memory      = "256M"
    max_instance_count    = 10
    timeout_seconds       = 30
    service_account_email = google_service_account.function_status.email
    environment_variables = {
      GCP_PROJECT_ID           = "${google_project.project.name}",
      GCP_ZONE                 = "${var.default_zone}",
      GCP_INSTANCE_NAME        = "${var.server_name}",
      USER_NAME                = "guest",
      PASSWORD_SEC_NAME        = "server-manager-password",
      TELNET_HOST              = "${google_compute_address.static_external_ip.address}",
      TELNET_PORT              = "8081"
      TELNET_PASSWORD_SEC_NAME = "telnet-password"
    }
  }
}

resource "google_cloud_run_service_iam_binding" "function_status_public_access" {
  project  = google_project.project.name
  location = google_cloudfunctions2_function.function_status.location
  service  = google_cloudfunctions2_function.function_status.name
  role     = "roles/run.invoker"
  members = [
    "allUsers"
  ]
}

resource "google_cloudfunctions2_function_iam_member" "function_status_public_access" {
  project        = google_project.project.name
  location       = google_cloudfunctions2_function.function_status.location
  cloud_function = google_cloudfunctions2_function.function_status.name

  role   = "roles/cloudfunctions.invoker"
  member = "allUsers"

  depends_on = [google_cloudfunctions2_function.function_status]
}

resource "google_service_account" "function_status" {
  project    = google_project.project.name
  account_id = "cloud-function-status"
}

resource "google_project_iam_member" "function_status" {
  for_each = toset(local.function_status_roles)
  project  = google_project.project.name
  member   = "serviceAccount:${google_service_account.function_status.email}"
  role     = each.value
}

data "archive_file" "function_status_source" {
  type        = "zip"
  output_path = "upload/function-status-source.zip"
  source_dir  = "functions/status/"
  excludes    = ["functions/status/cmd/main.go", "functions/status/cmd/README.md"]
}

resource "google_storage_bucket_object" "function_status_source" {
  name   = "function-status-source.zip"
  bucket = google_storage_bucket.cloud_function_staging_bucket.name
  source = "upload/function-status-source.zip" # Add path to the zipped function source code

  depends_on = [data.archive_file.function_status_source]
}
