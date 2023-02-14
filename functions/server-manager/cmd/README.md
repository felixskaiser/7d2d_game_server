# Run locally

>Needs GCP Application Default Credentials setup for project `game-server-7d2d-felix`.

Set required env vars:

```sh
export GCP_PROJECT_ID=game-server-7d2d-felix
export GCP_ZONE=europe-west3-b
export GCP_INSTANCE_NAME=game-server-7d2d
export USER_NAME=guest
export PASSWORD_SEC_NAME=server-manager-password
export TELNET_HOST=35.242.228.181
export TELNET_PORT=8081
export TELNET_PASSWORD_SEC_NAME=telnet-password
```

Run server:

```sh
go run cmd/main.go
```
