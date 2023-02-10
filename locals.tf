locals {
  # network
  public_ip_range          = "0.0.0.0/0"
  lowest_fw_rule_priority  = "65533"
  default_fw_rule_priority = "10000"
  game_server_port         = "26900"
  game_server_ports        = ["26900-26903"]
  game_server_network_tags = ["game-server"]

  # game server
  game_server_startup_script = templatefile("./server_startup_script.tftpl", {
    SERVERCONFIG_DEFAULT  = local.serverconfig_game_default,
    SERVERCONFIG_OFFHOURS = local.serverconfig_game_offhours,
    SERVER_CMD_SCRIPT     = base64encode(file("./cmd/7d2d_server")),
    ADMINCONFIG           = base64encode(file("./config/serveradmin.xml"))
    }
  )

  serverconfig_game_default = base64encode(templatefile("./config/serverconfig_game_default.xml.tftpl",
    {
      SERVERCONFIG_BASE = local.serverconfig_base
    }
  ))

  serverconfig_game_offhours = base64encode(templatefile("./config/serverconfig_game_offhours.xml.tftpl",
    {
      SERVERCONFIG_BASE = local.serverconfig_base
    }
  ))

  serverconfig_base = templatefile("./config/serverconfig_base.xml.tftpl",
    {
      SERVER_PASSWORD        = random_password.server_password.result
      CONTROL_PANEL_PASSWORD = random_password.control_panel_password.result
    }
  )

  game_server_shutdown_script = "#! /bin/bash 7d2d_server stop"
}
