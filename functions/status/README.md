# Game Server Status Cloud Function

A simple GCP Cloud Function to check if the game server is running and if players are online without having to join the game.

Deploy manually:

```sh
gcloud functions deploy status-test \
--project=game-server-7d2d-felix \
--region=europe-west3 \
--gen2 \
--stage-bucket=game-server-7d2d-felix-cloud-functions-stage \
--service-account=cloud-function-status@game-server-7d2d-felix.iam.gserviceaccount.com \
--runtime=go119 \
--set-env-vars=GCP_PROJECT_ID=game-server-7d2d-felix,USER_NAME=guest,PASSWORD_SEC_NAME=server-manager-password,GCP_ZONE=europe-west3-b,GCP_INSTANCE_NAME=game-server-7d2d \
--source=. \
--entry-point=Entrypoint \
--allow-unauthenticated \
--trigger-http
```
