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
