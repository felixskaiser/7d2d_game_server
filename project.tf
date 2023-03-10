resource "google_project" "project" {
  name                = var.project_id
  project_id          = var.project_id
  billing_account     = data.google_billing_account.billing.id
  auto_create_network = false

  lifecycle {
    prevent_destroy = true
  }
}

data "google_billing_account" "billing" {
  billing_account = var.billing_account_id
}

resource "google_project_service" "secret_manager_service" {
  service    = "secretmanager.googleapis.com"
  project    = google_project.project.name
  depends_on = [google_project.project]
}

resource "google_project_service" "cloud_functions_service" {
  service    = "cloudfunctions.googleapis.com"
  project    = google_project.project.name
  depends_on = [google_project.project]
}

resource "google_project_service" "cloud_build_service" {
  service    = "cloudbuild.googleapis.com"
  project    = google_project.project.name
  depends_on = [google_project.project]
}

resource "google_project_service" "artifact_registry_service" {
  service    = "artifactregistry.googleapis.com"
  project    = google_project.project.name
  depends_on = [google_project.project]
}

resource "google_project_service" "cloud_run_service" {
  service    = "run.googleapis.com"
  project    = google_project.project.name
  depends_on = [google_project.project]
}

resource "google_project_service" "monitoring_service" {
  service    = "monitoring.googleapis.com"
  project    = google_project.project.name
  depends_on = [google_project.project]
}
