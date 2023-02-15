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
# Server Manager Cloud Function
###############################################################################

resource "google_cloudfunctions2_function" "server_manager_function" {
  name     = "server-7daystodie"
  project  = google_project.project.name
  location = var.default_region

  build_config {
    runtime     = "go119"
    entry_point = "Entrypoint"

    source {
      storage_source {
        bucket = google_storage_bucket.cloud_function_staging_bucket.name
        object = google_storage_bucket_object.server_manager_function_source.name
      }
    }

    # Update Cloud Function if source code changes
    environment_variables = {
      SOURCE_HASH = "${google_storage_bucket_object.server_manager_function_source.md5hash}"
    }
  }

  service_config {
    available_memory      = "256M"
    max_instance_count    = 10
    timeout_seconds       = 30
    service_account_email = google_service_account.server_manager_function.email
    environment_variables = {
      IS_CLOUD_FUNCTION        = "true"
      GCP_PROJECT_ID           = "${google_project.project.name}",
      GCP_ZONE                 = "${var.default_zone}",
      GCP_INSTANCE_NAME        = "${var.server_name}",
      USER_NAME                = "admin",
      PASSWORD_SEC_NAME        = "server-manager-password",
      TELNET_HOST              = "${google_compute_address.static_external_ip.address}",
      TELNET_PORT              = "8081"
      TELNET_PASSWORD_SEC_NAME = "telnet-password"
    }
  }

  depends_on = [google_storage_bucket_object.server_manager_function_source]
}

resource "google_cloud_run_service_iam_member" "server_manager_function_public_access" {
  project  = google_project.project.name
  location = google_cloudfunctions2_function.server_manager_function.location
  service  = google_cloudfunctions2_function.server_manager_function.name
  role     = "roles/run.invoker"
  member   = "allUsers"

  depends_on = [google_cloudfunctions2_function.server_manager_function]
}

resource "google_cloudfunctions2_function_iam_member" "server_manager_function_public_access" {
  project        = google_project.project.name
  location       = google_cloudfunctions2_function.server_manager_function.location
  cloud_function = google_cloudfunctions2_function.server_manager_function.name
  role   = "roles/cloudfunctions.invoker"
  member = "allUsers"

  depends_on = [google_cloudfunctions2_function.server_manager_function]
}

resource "google_service_account" "server_manager_function" {
  project    = google_project.project.name
  account_id = "server-manager-function"
}

resource "google_project_iam_member" "server_manager_function" {
  for_each = toset(local.server_manager_function_roles)
  project  = google_project.project.name
  member   = "serviceAccount:${google_service_account.server_manager_function.email}"
  role     = each.value
}

resource "google_storage_bucket_object" "server_manager_function_source" {
  name   = "server-manager-function-source.zip"
  bucket = google_storage_bucket.cloud_function_staging_bucket.name
  source = "upload/server-manager-function-source.zip" # Add path to the zipped function source code

  depends_on = [data.archive_file.server_manager_function_source]
}

data "archive_file" "server_manager_function_source" {
  type        = "zip"
  output_path = "upload/server-manager-function-source.zip"
  source_dir  = "functions/server-manager/"
  excludes    = ["functions/server-manager/cmd/main.go", "functions/server-manager/cmd/README.md"]
}
