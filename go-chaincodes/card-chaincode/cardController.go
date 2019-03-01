/*
==== Install/Instantiate/Upgrade
peer chaincode install -n cc-card -p github.com/go-chaincodes/card-chaincode -v v1
peer chaincode instantiate -o orderer.example.com:7050 -C mychannel -n cc-card -c '{"Args":["init"]}' -v v1
peer chaincode upgrade -o orderer.example.com:7050 -C mychannel -n cc-card -c '{"Args":["init"]}' -v v2


==== Cards ====
peer chaincode invoke -C mychannel -n cc-card -c '{"Args":["CreateCard","10","1"]}'

peer chaincode query -C mychannel -n cc-card -c '{"Args":["GetCardByNumber","10"]}'
*/
package main

import (
	"fmt"

	card "github.com/go-chaincodes/card-chaincode/card"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// CardChaincode struct
type CardChaincode struct {
}

/*
 * ============================================================
 *  Main
 * ============================================================
 */
func main() {
	err := shim.Start(new(CardChaincode))
	if err != nil {
		fmt.Printf("Erro ao iniciar Card Chaincode: %s", err)
	}
}

/*
 * ============================================================
 * Init - initializes chaincode
 * ============================================================
 */
func (t *CardChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

/*
 * ============================================================
 * Invoke - Entry point for Invocations
 * ============================================================
 */
func (t *CardChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("Card Invoke is running " + function)

	// Handle different functions
	switch function {
	case "CreateCard":
		return card.CreateCard(stub, args)
	case "GetCardByNumber":
		return card.GetCardByNumber(stub, args)
	default:
		// error
		return shim.Error("Received unknown function invocation on Card Chaincode")
	}
}
