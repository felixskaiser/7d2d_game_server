output "function_status_uri" {
  value = google_cloudfunctions2_function.function_status.service_config[0].uri
}
