# 7d2d game server setup on GCP

```sh
terraform apply -var-file=config/default.tfvars -var 'billing_account_id=<my-gcp-billing-account-id>'
```

## Start game server

```sh
su steam
#TODO: patch serverconfig.xml
cd /home/steam/7d2d
LOGFILE=$(pwd)/7DaysToDieServer_Data/output_log__$(date +%Y-%m-%d__%H-%M-%S).txt
./7DaysToDieServer.x86_64 -logfile $LOGFILE -quit -batchmode -nographics -dedicated -configfile=serverconfig.xml
```

## Check logs

```sh
ls /home/steam/7d2d/7DaysToDieServer_Data/|grep "output_log"
cat <file>
```
