#!/bin/bash
#
# License: Apache-2.0
#

# Exit on first error, print all commands.
set -ev

docker-compose -f docker-compose.yml down
docker-compose -f docker-compose.yml up -d

# wait for Hyperledger Fabric to start
# incase of errors when running later commands, issue export FABRIC_START_TIMEOUT=<larger number>
export FABRIC_START_TIMEOUT=6
sleep ${FABRIC_START_TIMEOUT}

# Create the channel
docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@org1.example.com/msp" peer0.org1.example.com peer channel create -o orderer.example.com:7050 -c mychannel -f /etc/hyperledger/configtx/channel.tx

# Join peer0.org1.example.com to the channel.
docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@org1.example.com/msp" peer0.org1.example.com peer channel join -b mychannel.block

# Splunk
# docker run --rm --name splunk -d -p 8000:8000 -p 5514:5514 -e 'SPLUNK_START_ARGS=--accept-license' -e 'SPLUNK_PASSWORD=admin123' splunk/splunk:latest

# Logspout
# docker run --rm --name logspout -d -v /var/run/docker.sock:/var/run/docker.sock -p 127.0.0.1:4000:80 gliderlabs/logspout:latest syslog+tcp://127.0.0.1:9514
