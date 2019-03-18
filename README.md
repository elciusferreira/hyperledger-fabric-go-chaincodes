# Hyperledger Fabric Chaincode with Go

This repository presents three Hyperledger Fabric (v1.4) Chaincodes (SmartContracts) written using Go: account, card and transfer.

The account chaincode allows a simple account creation and query. The card chaincode allows a simple card creation (to an existent account) and query. The transfer chaincode allows money transfer from one account to another.

- - -

## Setting up the environment and deploying the network

Install Hyperledger Fabric prerequisites (including Go):

<https://hyperledger-fabric.readthedocs.io/en/latest/prereqs.html>

Install Fabric:

<https://hyperledger-fabric.readthedocs.io/en/latest/install.html>

Open a terminal and run the following commands to create the necessary directories:

    mkdir -p ~/go/src/github.com/hyperledger
    cd ~/go/src/github.com/hyperledger
    pwd

You should see something like the following path:

    /home/local/your_username/go/src/github.com/hyperledger

Then, clone hyperledger fabric repo inside that path:

    git clone https://github.com/hyperledger/fabric.git

Return to `github.com` folder and clone this repo:

    cd ~/go/src/github.com
    git clone https://github.com/elciusferreira/hyperledger-fabric-go-chaincodes.git

Everything should be good to go from now :)

### Note about GOPATH

To things work properly, we supposed your `GOPATH` environment variable is set to `/home/local/your_username/go`. To make sure it is, run:

    echo $GOPATH

It should print your configured go path, which go will look up for source code of your imports. That's why go returns an error while trying to build if you don't set it properly, and it cannot find fabric dependencies.

If your `GOPATH` is set differently, configure your `github.com` source directories accordingly.

- - -

## Starting the network

Go to `hyperledger-fabric-go-chaincodes/basic-network/` directory and start the network containers by running the `start.sh` script:

    cd $GOPATH/src/github.com/hyperledger-fabric-go-chaincodes/basic-network/
    ./start.sh

Last line of output should be:

    executeJoin -> INFO 002 Successfully submitted proposal to join channel

This basic-network is a simple infrastructure that consistis of one peer and one orderer. You can check the network containers (peer, ca, orderer, cli and couchdb) by running the docker command:

    docker ps

- - -

## Installing and instantiating chaincodes

To execute the following commands, you should enter the facbric cli container. To do so, run:

    docker exec -it cli bash

You should see the following:

    root@xxxxxxxxxxxx:/opt/gopath/src/github.com/hyperledger/fabric/peer#

Where `xxxxxxxxxxxx` is the hash id of the cli container.
Now, you're able to use the cli to install and instantiate the chaincodes:

    peer chaincode install -n cc-account -p github.com/hyperledger-fabric-go-chaincodes/account-chaincode -v v1
    peer chaincode instantiate -o orderer.example.com:7050 -C mychannel -n cc-account -c '{"Args":["init"]}' -v v1

    peer chaincode install -n cc-card -p github.com/hyperledger-fabric-go-chaincodes/card-chaincode -v v1
    peer chaincode instantiate -o orderer.example.com:7050 -C mychannel -n cc-card -c '{"Args":["init"]}' -v v1

    peer chaincode install -n cc-transfer -p github.com/hyperledger-fabric-go-chaincodes/transfer-chaincode -v v1
    peer chaincode instantiate -o orderer.example.com:7050 -C mychannel -n cc-transfer -c '{"Args":["init"]}' -v v1

The peer chaincode install command sends the chaincode to the network peer. The peer chaincode instantiate command will build the go files and if there are no errors, the chaincode will be ready for use.

- - -

## Calling Chaincodes

Now that everything is configured and the chaincodes are installed and instantiated, you can call the chaincodes to execute operations.

### Account chaincode

With the account chaincode installed and instantiated you can create an account:

    peer chaincode invoke -C mychannel -n cc-account -c '{"Args":["Create","1","1000","Elcius"]}'

Where the first argument is the function name, the second is the unique account number, the third is the initial account balance and the last one is the account owner name.  

Create a predefined set of accounts:

    peer chaincode invoke -C mychannel -n cc-account -c '{"Args":["Init"]}'

Query an account by its number:

    peer chaincode query -C mychannel -n cc-account -c '{"Args":["GetByNumber","1"]}'

Celete an account by its number:

    peer chaincode invoke -C mychannel -n cc-account -c '{"Args":["Delete","1"]}'

Get a history for an account by its number:

    peer chaincode invoke -C mychannel -n cc-account -c '{"Args":["GetHistory","1"]}'

Get an account by owner name:

    peer chaincode query -C mychannel -n cc-account -c '{"Args":["GetByOwner","Elcius"]}'

### Card chaincode

With the Card chaincode installed and instantiated you can create a card:

    peer chaincode invoke -C mychannel -n cc-card -c '{"Args":["Create","10","1"]}'

Where the first argument is the function name, the second is the card number and the last one is the existent account number related to the card to be created.

Query a card by its number:

    peer chaincode query -C mychannel -n cc-card -c '{"Args":["GetByNumber","10"]}'

### Transfer chaincode

With the Transfer chaincode installed and instantiated you can transfer money from one account to another:

    peer chaincode invoke -C mychannel -n cc-transfer -c '{"Args":["Money","1","2","500"]}'

Where the first argument is the function name, the second is the payer account number, the second is the receiver account number and the last one is the money amount to be transfered.

- - -

## Other instructions

If you want to edit the code and test the changes, you should build your go files. To do so, make sure your `GOPATH`, `GOROOT` and `PATH` environment variables are properly set in `.bashrc` file. Check your `.bashrc` file in your home directory (`/home/local/your_username/`):

    nano ~/.bashrc

The following lines must be somewhere in the file:

    export GOPATH=$HOME/go/
    export GOROOT=/usr/local/go
    export PATH=$PATH:$GOROOT/bin

To build the code and check if there are any errors, go to the modified chaincode directory and run:

    cd $GOPATH/src/github.com/hyperledger-fabric-go-chaincodes/account-chaincode/
    go build

If there are no errors you can proceed and use the cli again to install the edited chaincode on the peer and upgrade the network with your new chaincode version. For example, if the Account chaincode is modified (run inside cli container):

    peer chaincode install -n cc-account -p github.com/hyperledger-fabric-go-chaincodes/account-chaincode -v v2
    peer chaincode upgrade -o orderer.example.com:7050 -C mychannel -n cc-account -c '{"Args":["init"]}' -v v2

To see the chaincode logs, run:

    docker logs -f <dev_container_name>

To check the dev container name of each chaincode installed, run the following docker command:

    docker ps -f name=dev-*

For example, to see the logs of account chaincode:

    docker logs -f dev-peer0.org1.example.com-cc-account-v1

To shutdown the network completely, go to the fabric-samples/basic-network directory and run:

    ./stop.sh
    ./teardown.sh
