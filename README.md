# 7d2d game server setup on GCP

## Setup infrastructure

Configure files in `config/` as required:

- `default.tfvars`: Terraform variables used GCP infrastructure setup
- `serverconfig_base.xml.tftpl`: Base game server config (check <https://7daystodie.fandom.com/wiki/Server>), used for generating final `serverconfig_game_<mode>.xml` files
- `serverconfig_game_default.xml.tftpl`: Default gameplay setting server config (check <https://7daystodie.fandom.com/wiki/Server>), will be combined with `serverconfig_base.xml.tftpl` for generating final `serverconfig_game_default.xml` file
- `serverconfig_game_offhours.xml.tftpl`: Slower gameplay setting server config (check <https://7daystodie.fandom.com/wiki/Server>), will be combined with `serverconfig_base.xml.tftpl` for generating final `serverconfig_game_offhours.xml` file

```sh
terraform apply -var-file=config/default.tfvars -var 'billing_account_id=<my-gcp-billing-account-id>'
```

## Start game server

SSH into game server, then:

```sh
7d2d_server start serverconfig_game_default.xml
```

## Stop game server

SSH into game server, then:

```sh
7d2d_server stop
```

## Check logs

SSH into game server, then find logs at `/home/steam/7d2d/7DaysToDieServer_Data/logs`.
