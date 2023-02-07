###############################################################################
# Network
###############################################################################

resource "google_compute_network" "network" {
  project                 = google_project.project.name
  name                    = "game-server-network"
  auto_create_subnetworks = false
  mtu                     = 1460
}

resource "google_compute_subnetwork" "subnet" {
  project       = google_project.project.name
  name          = "game-server-subnet"
  ip_cidr_range = "10.156.0.0/20"
  region        = var.default_region
  network       = google_compute_network.network.id
}

resource "google_compute_address" "static_external_ip" {
  project = google_project.project.name
  name    = "game-server-address"
  region  = var.default_region
}

###############################################################################
# Firewall - project defaults
###############################################################################

resource "google_compute_firewall" "ingress_deny_all_default" {
  project       = google_project.project.name
  name          = "ingress-deny-all-default"
  network       = google_compute_network.network.self_link
  direction     = "INGRESS"
  source_ranges = [local.public_ip_range]
  priority      = local.lowest_fw_rule_priority

  deny {
    protocol = "all"
    # no specifying ports means all ports
  }
}

resource "google_compute_firewall" "egress_deny_all_default" {
  project            = google_project.project.name
  name               = "egress-deny-all-default"
  network            = google_compute_network.network.self_link
  direction          = "EGRESS"
  destination_ranges = [local.public_ip_range]
  priority           = local.lowest_fw_rule_priority

  deny {
    protocol = "all"
    # no specifying ports means all ports
  }
}

resource "google_compute_firewall" "ingress_allow_standard_internal" {
  project       = google_project.project.name
  name          = "ingress-allow-standard-internal"
  network       = google_compute_network.network.self_link
  direction     = "INGRESS"
  source_ranges = [local.public_ip_range]
  priority      = local.default_fw_rule_priority

  allow {
    protocol = "tcp"
    ports    = ["0-65535"]
  }

  allow {
    protocol = "udp"
    ports    = ["0-65535"]
  }

  allow {
    protocol = "icmp"
    # no specifying ports means all ports
  }
}

###############################################################################
# Firewall - game server
###############################################################################

resource "google_compute_firewall" "egress_allow_standard_public" {
  project            = google_project.project.name
  name               = "egress-allow-standard-public"
  network            = google_compute_network.network.self_link
  direction          = "EGRESS"
  destination_ranges = [local.public_ip_range]
  priority           = local.default_fw_rule_priority
  target_tags        = local.game_server_network_tags

  allow {
    protocol = "tcp"
    ports    = ["80", "443"]
  }

  allow {
    protocol = "icmp"
    # no specifying ports means all ports
  }
}

resource "google_compute_firewall" "ingress_allow_icmp_public" {
  project       = google_project.project.name
  name          = "ingress-allow-icmp-public"
  network       = google_compute_network.network.self_link
  direction     = "INGRESS"
  source_ranges = [local.public_ip_range]
  priority      = local.default_fw_rule_priority
  target_tags   = local.game_server_network_tags

  allow {
    protocol = "icmp"
    # no specifying ports means all ports
  }
}

resource "google_compute_firewall" "ingress_allow_game_server_public" {
  project       = google_project.project.name
  name          = "ingress-allow-game-server-public"
  network       = google_compute_network.network.self_link
  direction     = "INGRESS"
  source_ranges = [local.public_ip_range]
  priority      = local.default_fw_rule_priority
  target_tags   = local.game_server_network_tags

  allow {
    protocol = "tcp"
    ports    = [local.game_server_port]
  }

  allow {
    protocol = "udp"
    ports    = local.game_server_ports
  }
}

#TODO: test if required
resource "google_compute_firewall" "egress_allow_game_server_public" {
  project            = google_project.project.name
  name               = "egress-allow-game-server-public"
  network            = google_compute_network.network.self_link
  direction          = "EGRESS"
  destination_ranges = [local.public_ip_range]
  priority           = local.default_fw_rule_priority
  target_tags        = local.game_server_network_tags

  allow {
    protocol = "tcp"
    ports    = [local.game_server_port]
  }

  allow {
    protocol = "udp"
    ports    = local.game_server_ports
  }
}

resource "google_compute_firewall" "ingress_allow_control_plane_public" {
  project       = google_project.project.name
  name          = "ingress-allow-control-plane-public"
  network       = google_compute_network.network.self_link
  direction     = "INGRESS"
  source_ranges = [local.public_ip_range]
  priority      = local.default_fw_rule_priority
  target_tags   = local.game_server_network_tags

  allow {
    protocol = "tcp"
    ports    = ["8080"]
  }
}
