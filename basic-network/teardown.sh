#!/bin/bash
#
# License: Apache-2.0
# 

# Exit on first error, print all commands.
set -ev

# Shut down the Docker containers for the system tests.
docker-compose -f docker-compose.yml kill && docker-compose -f docker-compose.yml down || true

# kill and remove all containers
docker kill $(docker ps -aq) || true
docker rm $(docker ps -aq) || true

# remove all volumes
docker volume rm $(docker volume ls -q) || true

# remove chaincode docker images
docker rmi $(docker images dev-* -q) || true
