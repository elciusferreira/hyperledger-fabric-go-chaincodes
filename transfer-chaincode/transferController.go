/*
==== Install/Instantiate/Upgrade
peer chaincode install -n cc-transfer -p github.com/hyperledger-fabric-go-chaincodes/transfer-chaincode -v v1
peer chaincode instantiate -o orderer.example.com:7050 -C mychannel -n cc-transfer -c '{"Args":["init"]}' -v v1
peer chaincode upgrade -o orderer.example.com:7050 -C mychannel -n cc-transfer -c '{"Args":["init"]}' -v v2

-- Invoke
peer chaincode invoke -C mychannel -n cc-transfer -c '{"Args":["Money","1","2","500"]}'

*/
package main

import (
	"fmt"

	"github.com/hyperledger-fabric-go-chaincodes/transfer-chaincode/transfer"

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
		fmt.Printf("Erro ao iniciar Transfer Chaincode: %s", err)
	}
}

// Init - initializes chaincode
func (t *TransferController) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

// Invoke - Entry point for Invocations
func (t *TransferController) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("Transfer Invoke is running " + function)

	// Handle different functions
	switch function {
	case "Money":
		return transfer.Money(stub, args)
	default:
		// error
		return shim.Error("Received unknown function invocation on Transfer Chaincode")
	}
}
