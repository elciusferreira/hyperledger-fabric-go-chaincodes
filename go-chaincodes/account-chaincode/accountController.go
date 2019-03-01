/*
==== Install/Instantiate/Upgrade
peer chaincode install -n cc-account -p github.com/go-chaincodes/account-chaincode -v v1
peer chaincode instantiate -o orderer.example.com:7050 -C mychannel -n cc-account -c '{"Args":["init"]}' -v v1
peer chaincode upgrade -o orderer.example.com:7050 -C mychannel -n cc-account -c '{"Args":["init"]}' -v v1


==== Accounts ====
peer chaincode invoke -C mychannel -n cc-account -c '{"Args":["InitAccounts"]}'
peer chaincode invoke -C mychannel -n cc-account -c '{"Args":["CreateAccount","1","1000","Elcius"]}'
peer chaincode invoke -C mychannel -n cc-account -c '{"Args":["CreateAccount","2","1000","Natan"]}'

peer chaincode query -C mychannel -n cc-account -c '{"Args":["GetAccountByNumber","1"]}'
*/

package main

import (
	"fmt"

	acc "github.com/go-chaincodes/account-chaincode/account"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// AccountsChaincode struct
type AccountsChaincode struct {
}

/*
 * ============================================================
 *  Main
 * ============================================================
 */
func main() {
	err := shim.Start(new(AccountsChaincode))
	if err != nil {
		fmt.Printf("Error initializing Accounts Chaincode: %s", err)
	}
}

/*
 * ============================================================
 * Init - initializes chaincode
 * ============================================================
 */
func (t *AccountsChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

/*
 * ============================================================
 * Invoke - Entry point for Invocations
 * ============================================================
 */
func (t *AccountsChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("Accounts Invoke is running " + function)

	// Handle different functions
	switch function {
	case "InitAccounts":
		return acc.InitAccounts(stub)
	case "CreateAccount":
		return acc.CreateAccount(stub, args)
	case "GetAccountByNumber":
		return acc.GetAccountByNumber(stub, args)
	case "UpdateAccountByNumber":
		return acc.UpdateAccountByNumber(stub, args)
	default:
		// error
		return shim.Error("Received unknown function invocation on Account Chaincode")
	}
}
