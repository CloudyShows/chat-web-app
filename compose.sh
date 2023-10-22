#!/bin/bash
DOCKER_CONFIG="/usr/local/lib/docker"
sudo mkdir -p "$DOCKER_CONFIG/cli-plugins"
sudo curl -SL https://github.com/docker/compose/releases/download/v2.23.0/docker-compose-linux-x86_64 -o "$DOCKER_CONFIG"/cli-plugins/docker-compose
sudo chmod +x /usr/local/lib/docker/cli-plugins/docker-compose
docker compose version