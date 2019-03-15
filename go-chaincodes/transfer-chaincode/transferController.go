/*
==== Install/Instantiate/Upgrade
peer chaincode install -n cc-transfer -p github.com/go-chaincodes/transfer-chaincode -v v1.0.0
peer chaincode instantiate -o orderer.example.com:7050 -C mychannel -n cc-transfer -c '{"Args":["init"]}' -v v1.0.0
peer chaincode upgrade -o orderer.example.com:7050 -C mychannel -n cc-transfer -c '{"Args":["init"]}' -v v1.0.1

==== List chaincodes ====
peer chaincode list --installed
peer chaincode list --instantiated -C mychannel

==== Transfer ====
 +++ Invokes
peer chaincode invoke -C mychannel -n cc-transfer -c '{"Args":["Money","1","2","500"]}'
*/

package main

import (
	"fmt"

	"github.com/go-chaincodes/transfer-chaincode/transfer"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// TransferController struct
type TransferController struct {
}

//  Main
func main() {
	err := shim.Start(new(TransferController))
	if err != nil {
		fmt.Println("failed to initialize transfer chaincode" + err.Error())
	}
}

// Init - initializes chaincode
func (t *TransferController) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

// Invoke - Entry point for Invocations
func (t *TransferController) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("[DEBUG] Transfer chaincode invoking " + function + " function")

	// Handle different functions
	switch function {
	case "Money":
		return transfer.Money(stub, args)
	default:
		return shim.Error("received unknown function invocation on transfer chaincode")
	}
}
