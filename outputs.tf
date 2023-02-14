output "function_status_uri" {
  value = google_cloudfunctions2_function.server_manager_function.service_config[0].uri
}
