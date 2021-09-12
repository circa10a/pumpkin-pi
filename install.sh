#!/usr/bin/env bash

# Install docker if not installed
if [ ! -x "$(command -v docker)" ]; then
    curl -fsSL https://get.docker.com -o get-docker.sh
    sudo sh get-docker.sh
fi
# Add current user to docker group
sudo usermod -aG docker ${USER}
# Install docker-compose
sudo apt-get install -y libffi-dev libssl-dev
sudo apt install -y python3-dev
sudo apt-get install -y python3 python3-pip
sudo pip3 install docker-compose
# Install git
sudo apt install git -y
# Clone repository
git clone https://github.com/circa10a/pumpkin-pi.git
cd pumpkin-pi/
# Start pumpkin-pi + pi-blaster
docker-compose up -d