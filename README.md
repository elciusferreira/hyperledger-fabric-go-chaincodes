# Hyperledger Fabric Chaincode with Go

This repository presents three Hyperledger Fabric (v1.4) Chaincodes (SmartContracts) written using Go: Account, card and transfer.

The account chaincode allows a simple account creation and query. The card chaincode allows a simple card creation (to an existent account) and query. The transfer chaincode allows money transfer from one account to another.

## Setting up the environment and deploying the network

Hyperledger Fabric prerequisites:

<https://hyperledger-fabric.readthedocs.io/en/latest/prereqs.html>

Install Fabric:

<https://hyperledger-fabric.readthedocs.io/en/latest/install.html>

Clone fabric-samples repo:

    git clone https://github.com/hyperledger/fabric-samples.git

Clone this repo:

    git clone https://github.com/elciusferreira/hyperledger-fabric-go-chaincodes.git

Copy go-chaincodes directory to fabric-samples/chaincode.
Inside fabric-samples/basic-network folder edit start.sh file. In the file, go to the command:

    docker-compose -f docker-compose.yml up -d ca.example.com orderer.example.com peer0.org1.example.com couchdb
and remove this part:

    ca.example.com orderer.example.com peer0.org1.example.com couchdb
So, the command must be only:

    docker-compose -f docker-compose.yml up -d
Save and close start.sh.

Navegate using the terminal to fabric-samples/basic-network directory and type:

    chmod +x start.sh stop.sh teardown.sh

Then, start the network containers:

    ./start.sh
   Check them by typing the docker command:

    docker ps

Enter in cli container:

    docker exec -it cli bash

Use the cli to install and instantiate the chaincodes:

    peer chaincode install -n cc-account -p github.com/go-chaincodes/account-chaincode -v v1
    peer chaincode instantiate -o orderer.example.com:7050 -C mychannel -n cc-account -c '{"Args":["init"]}' -v v1

    peer chaincode install -n cc-card -p github.com/go-chaincodes/card-chaincode -v v1
    peer chaincode instantiate -o orderer.example.com:7050 -C mychannel -n cc-card -c '{"Args":["init"]}' -v v1

    peer chaincode install -n cc-transfer -p github.com/go-chaincodes/transfer-chaincode -v v1
    peer chaincode instantiate -o orderer.example.com:7050 -C mychannel -n cc-transfer -c '{"Args":["init"]}' -v v1

## Account chaincode

To create an account:

    peer chaincode invoke -C mychannel -n cc-account -c '{"Args":["CreateAccount","1","1000","Elcius"]}'

Where the first argument is the function name, the second is the unique account number, the third is the initial account balance and the last one is the account owner name.  

To create a predefined set of accounts:

    peer chaincode invoke -C mychannel -n cc-account -c '{"Args":["Init"]}'

To query an account by its number:

    peer chaincode query -C mychannel -n cc-account -c '{"Args":["GetByNumber","1"]}'

To delete an account by its number:

    peer chaincode invoke -C mychannel -n cc-account -c '{"Args":["Delete","1"]}'

To get a history for an account by its number:

    peer chaincode invoke -C mychannel -n cc-account -c '{"Args":["GetHistory","1"]}'

## Card chaincode

To create a card:

    peer chaincode invoke -C mychannel -n cc-card -c '{"Args":["Create","10","1"]}'

Where the first argument is the function name, the second is the card number and the last one is the existent account number related to the new card.

To query a card by its number:

    peer chaincode query -C mychannel -n cc-card -c '{"Args":["GetByNumber","10"]}'

## Transfer chaincode

To transfer money from one account to another:

    peer chaincode invoke -C mychannel -n cc-transfer -c '{"Args":["Money","1","2","500"]}'

Where the first argument is the function name, the second is the payer account number, the second is the receiver account number and the last one is the money amount to be transfered.

## Other instructions

If you want to edit the code and test, you should use the cli again to install on the peer the edited chaincode and upgrade the network with the chaincode new version number. For example, if the account chaincode is modified:

    peer chaincode install -n cc-account -p github.com/go-chaincodes/account-chaincode -v v2
    peer chaincode upgrade -o orderer.example.com:7050 -C mychannel -n cc-account -c '{"Args":["init"]}' -v v2

To see the log of a chaincode, open a new terminal tab/window and type the command:

    docker logs -f <dev_container_name>

To check the dev container name of each chaincode installed, type one more time the following docker command:

    docker ps

For example, to see the logs of account chaincode:

    docker logs -f dev-peer0.org1.example.com-cc-account-v1

To shutdown the network completely, go to the fabric-samples/basic-network directory and run:

    ./stop.sh
    ./teardown.sh
