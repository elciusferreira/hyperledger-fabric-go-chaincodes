/*
==== Install/Instantiate/Upgrade
peer chaincode install -n cc-account -p github.com/hyperledger-fabric-go-chaincodes/account-chaincode -v v1
peer chaincode instantiate -o orderer.example.com:7050 -C mychannel -n cc-account -c '{"Args":["init"]}' -v v1
peer chaincode upgrade -o orderer.example.com:7050 -C mychannel -n cc-account -c '{"Args":["init"]}' -v v2

==== List chaincodes ====
peer chaincode list --installed
peer chaincode list --instantiated -C mychannel

==== Accounts ====
 +++ Invokes
peer chaincode invoke -C mychannel -n cc-account -c '{"Args":["Init"]}'
peer chaincode invoke -C mychannel -n cc-account -c '{"Args":["Create","1","1000","Elcius"]}'
peer chaincode invoke -C mychannel -n cc-account -c '{"Args":["Create","2","1000","Natan"]}'
peer chaincode invoke -C mychannel -n cc-account -c '{"Args":["Create","6","1000","Marcelo"]}'
peer chaincode invoke -C mychannel -n cc-account -c '{"Args":["Delete","1"]}'
peer chaincode invoke -C mychannel -n cc-account -c '{"Args":["Update","{\"accountBalance\":7000,\"accountNumber\":2,\"accountOwner\":\"Natanael\",\"docType\":\"Account\"}"]}'

 +++ Queries
peer chaincode query -C mychannel -n cc-account -c '{"Args":["GetAll"]}' | jq
peer chaincode query -C mychannel -n cc-account -c '{"Args":["GetByNumber","1"]}' | jq
peer chaincode query -C mychannel -n cc-account -c '{"Args":["GetByOwner","Elcius"]}' | jq
peer chaincode query -C mychannel -n cc-account -c '{"Args":["GetHistory","1"]}' | jq
*/

package main

import (
	"fmt"

	"github.com/hyperledger-fabric-go-chaincodes/account-chaincode/account"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// AccountsChaincode struct
type AccountsChaincode struct {
}

// Main
func main() {
	err := shim.Start(new(AccountsChaincode))
	if err != nil {
		fmt.Println("failed to initialize accounts chaincode" + err.Error())
	}
}

// Init - initializes chaincode
func (t *AccountsChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

// Invoke - Entry point for Invocations
func (t *AccountsChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("[DEBUG] Accounts chaincode invoking " + function + " function")

	// Handle different functions
	switch function {
	case "Init":
		return account.Init(stub)
	case "Create":
		return account.Create(stub, args)
	case "GetAll":
		return account.GetAll(stub)
	case "GetByNumber":
		return account.GetByNumber(stub, args)
	case "GetByOwner":
		return account.GetByOwner(stub, args)
	case "Update":
		return account.Update(stub, args)
	case "Delete":
		return account.Delete(stub, args)
	case "GetHistory":
		return account.GetHistoryByAccNumber(stub, args)
	default:
		// Error
		return shim.Error("received unknown function invocation on account chaincode")
	}
}
