locals {
  # network
  public_ip_range          = "0.0.0.0/0"
  lowest_fw_rule_priority  = "65533"
  default_fw_rule_priority = "10000"
  game_server_port         = "26900"
  game_server_ports        = ["26900-26903"]
  game_server_network_tags = ["game-server"]
}
