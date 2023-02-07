#!/bin/bash
set -euxo pipefail

# IMPORTANT: Changing this script will force a recreation of the game server VM!

# Setup according to:
# - https://developer.valvesoftware.com/wiki/SteamCMD
# - https://7daystodie.fandom.com/wiki/Linux_Server_with_Amazon_Linux_AMI

###############################################################################
# Install 'steamcmd' dependencies
###############################################################################
echo "INFO: Install 'steamcmd' and dependencies"

# Add dependencies
sudo dpkg --add-architecture i386
sudo apt update
sudo apt install software-properties-common -y
sudo apt install lib32gcc-s1 -y
sudo apt install lib32stdc++6

###############################################################################
# Format game storage disk (if necessary) and mount it
###############################################################################
echo "INFO: Format game storage disk (if necessary) and mount it"

DEVICE_NAME="sdb"

# Only format disk if it's not formatted yet
# https://cloud.google.com/compute/docs/disks/add-persistent-disk#formatting
if sudo blkid --match-token TYPE=ext4 /dev/$DEVICE_NAME; then
    echo "disk '/dev/$DEVICE_NAME' is already formatted as ext4"
else
    echo "disk '$DEVICE_NAME' isn't formatted yet, formatting it as ext4..."
    sudo mkfs.ext4 -m 0 -E lazy_itable_init=0,lazy_journal_init=0,discard /dev/$DEVICE_NAME
fi


# Mount disk to /home/steam
# https://cloud.google.com/compute/docs/disks/add-persistent-disk#mounting
sudo mkdir -p /home/steam
sudo mount -o discard,defaults /dev/$DEVICE_NAME /home/steam

# Configure /etc/fstab
DISK_UUID=$(sudo blkid /dev/$DEVICE_NAME|sed 's~/dev/sdb: UUID="~~'|sed 's~" TYPE="ext4"~~')

if grep "$DISK_UUID" /etc/fstab; then
    echo "/etc/fstab already contains entry for disk '/dev/$DEVICE_NAME' with UUID '$DISK_UUID'"
else
    echo "/etc/fstab doesn't yet contain an entry for disk '/dev/$DEVICE_NAME' with UUID '$DISK_UUID', updating..."
    cp /etc/fstab fstab_new
    echo "UUID=$DISK_UUID /home/steam ext4 discard,defaults 0 2">>fstab_new
    sudo mv fstab_new /etc/fstab 
fi

# Grant write access to the disk for all users
sudo chmod a+w /home/steam

# Log results
sudo lsblk

###############################################################################
# Create and use user 'steam' to run 'steamcmd' with
###############################################################################
echo "INFO: Create and use user 'steam' to run 'steamcmd' with"

# Check if user steam exists, set it up if necessary
if [[ !$(cut -d: -f1 /etc/passwd|grep -q "steam") ]]; then
    echo "user 'steam' doesn't exist yet, setting it up..."
    # Add user and create home directory at /home/steam
    sudo useradd -m steam

    # Create user steam without password
    sudo passwd -d steam
else
    echo "user 'steam' already exists"
fi

# Switch to user steam and go to its home directory
su steam
cd /home/steam

###############################################################################
# Install 'steamcmd'
###############################################################################
echo "INFO: Check/install 'steamcmd'"

# IMPORTANT: Changing this directory may require manual cleanup
mkdir -p /home/steam/Steam
cd /home/steam/Steam

if [[ -e steamcmd.sh ]]; then
    echo "'steamcmd.sh' already exists on disk"
else
    echo "'steamcmd.sh' doesn't exist on disk yet, downloading..."
    # Download and extract SteamCMD for Linux
      curl -sqL "https://steamcdn-a.akamaihd.net/client/installer/steamcmd_linux.tar.gz" | tar zxvf -
fi

ls -alF
cat steamcmd.sh

###############################################################################
# Setup game server
###############################################################################
echo "INFO: Setup game server"

cd /home/steam/
# Ensure directories exist
# IMPORTANT: Changing these directories may require manual cleanup
mkdir -p /home/steam/7d2d
mkdir -p /home/steam/.steam/sdk64

# Link steamclient so that game server can find it
ln -s /home/steam/Steam/linux64/steamclient.so /home/steam/.steam/sdk64/steamclient.so

# Install/update 7D2D dedicated linux server at /home/steam/7d2d
# https://developer.valvesoftware.com/wiki/SteamCMD#Automating_SteamCMD
./Steam/steamcmd.sh +force_install_dir /home/steam/7d2d +login anonymous +app_update 294420 valdiate +quit

# Ensure user steam is owner of /home/steam 
sudo chown -R steam:steam /home/steam
