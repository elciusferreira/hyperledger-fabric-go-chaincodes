<h1>Hyperledger Fabric Chaincode with Go</h1>
This repository presents three Hyperledger Fabric (v1.4) Chaincodes (SmartContracts) written using Go: Account, card and transfer.

The account chaincode allows a simple account creation and query. The card chaincode allows a simple card creation (to an existent account) and query. The transfer chaincode allows money transfer from one account to another.

<h2>Setting up the environment</h2>
Install Hyperledger Fabric prerequisites (including Go): 

https://hyperledger-fabric.readthedocs.io/en/latest/prereqs.html

Install Fabric:

https://hyperledger-fabric.readthedocs.io/en/latest/install.html

Open a terminal and in root type the following commands to create the necessary directories:

    cd ~
    mkdir go && cd go
    mkdir src && cd src
    mkdir github.com && cd github.com
    mkdir hyperledger && cd hyperledger

Now, inside the brand new hyperledger folder type the command:
    
    pwd

The output should look like the path bellow:

    /home/local/your_user_name/go/src/github.com/hyperledger

Still inside hyperledger folder, clone fabric repo:

    ~/go/src/github.com/hyperledger$ git clone https://github.com/hyperledger/fabric.git

Return to github<!-- -->.com folder and clone my repo:

    ~/go/src/github.com/hyperledger$ cd ..
    ~/go/src/github.com$ git clone https://github.com/elciusferreira/hyperledger-fabric-go-chaincodes.git

<h2>Starting the network</h2>
Go to hyperledger-fabric-go-chaincodes/basic-network/ folder and start the network containers by running the start<!-- -->.sh scrypt:
	
    ~/go/src/github.com$ cd hyperledger-fabric-go-chaincodes/basic-network/
    ~/go/src/github.com/hyperledger/fabric-samples/basic-network$ ./start.sh

Last line of output should be:

    executeJoin -> INFO 002 Successfully submitted proposal to join channel

This basic-network is a simple infrastructure that consistis of one peer and one orderer. You can check the network containers (peer, ca, orderer, cli and couchdb) by typing the docker command:
	
    docker ps

<h2>Building and starting the chaincodes</h2>
Enter in cli container:
	
    docker exec -it cli bash

You should see the following:

    root@b067b942e2e5:/opt/gopath/src/github.com/hyperledger/fabric/peer#

Now, you are able to use the cli to install and instantiate the chaincodes:
	
    ...fabric/peer# peer chaincode install -n cc-account -p github.com/hyperledger-fabric-go-chaincodes/account-chaincode -v v1
    ...fabric/peer# peer chaincode instantiate -o orderer.example.com:7050 -C mychannel -n cc-account -c '{"Args":["init"]}' -v v1

	...fabric/peer# peer chaincode install -n cc-card -p github.com/hyperledger-fabric-go-chaincodes/card-chaincode -v v1
    ...fabric/peer# peer chaincode instantiate -o orderer.example.com:7050 -C mychannel -n cc-card -c '{"Args":["init"]}' -v v1

	...fabric/peer# peer chaincode install -n cc-transfer -p github.com/hyperledger-fabric-go-chaincodes/transfer-chaincode -v v1
    ...fabric/peer# peer chaincode instantiate -o orderer.example.com:7050 -C mychannel -n cc-transfer -c '{"Args":["init"]}' -v v1

The peer chaincode install command sends the chaincode to the network peer. The peer chaincode instantiate command will build the go files and if there are no errors, the chaincode will be ready for use.

<h2>Using Account chaincode</h2>
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



<h2>Using Card chaincode</h2>
With the Card chaincode installed and instantiated you can create a card:

	peer chaincode invoke -C mychannel -n cc-card -c '{"Args":["Create","10","1"]}'

Where the first argument is the function name, the second is the card number and the last one is the existent account number related to the card to be created.

Query a card by its number:

	peer chaincode query -C mychannel -n cc-card -c '{"Args":["GetByNumber","10"]}'

<h2>Using Transfer chaincode</h2>
With the Transfer chaincode installed and instantiated you can transfer money from one account to another:

	peer chaincode invoke -C mychannel -n cc-transfer -c '{"Args":["Money","1","2","500"]}'

Where the first argument is the function name, the second is the payer account number, the second is the receiver account number and the last one is the money amount to be transfered.

<h2> Other instructions </h2>
If you want to edit the code and test the changes, you should build your go files. To do that, make sure your GOPATH, GOROOT and PATH are properly set in .bashrc file. To check you can type from the root:

    ~$ nano ./bashrc

The following lines must be somewhere in the file:

    export GOPATH=/home/local/your_user_name/go/
    export GOROOT=/usr/local/go
    export PATH=$PATH:$GOROOT/bin

To know your_user_name you can type from the root on terminal:

    ~$ pwd

The output should be:

    /home/local/your_user_name

To build the code and check if there are any errors, navegate on terminal to the modified chaincode folder and type:

    ~$ cd go/src/github.com/hyperledger-fabric-go-chaincodes/account-chaincode/
    ~/go/src/github.com/hyperledger-fabric-go-chaincodes/account-chaincode$ go build

If there are no errors you can proceed and use the cli again to install the edited chaincode on the peer and upgrade the network with the chaincode new version number. For example, if the Account chaincode is modified:

	...fabric/peer# peer chaincode install -n cc-account -p github.com/hyperledger-fabric-go-chaincodes/account-chaincode -v v2
	...fabric/peer# peer chaincode upgrade -o orderer.example.com:7050 -C mychannel -n cc-account -c '{"Args":["init"]}' -v v2

To see the log of a chaincode (and all fmt.Println()), open a new terminal tab/window and type the command:

    docker logs -f <dev_container_name>

To check the dev container name of each chaincode installed, type one more time the following docker command:

    docker ps

For example, to see the logs of account chaincode:

    docker logs -f dev-peer0.org1.example.com-cc-account-v1

To shutdown the network completely, go to the fabric-samples/basic-network directory and run:

    ./stop.sh
    ./teardown.sh

