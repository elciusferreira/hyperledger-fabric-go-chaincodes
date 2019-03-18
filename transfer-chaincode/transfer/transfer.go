/*
Package transfer provides the money transfer service between accounts.
*/
package transfer

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger-fabric-go-chaincodes/account-chaincode/account"
	"github.com/hyperledger/fabric/common/util"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// Money - Transfer money between Accounts
// param: AccountNumber, AccountNumber, Value
func Money(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fmt.Println("[DEBUG] begin transfer.Money")

	var payerAcc account.Account
	var receiverAcc account.Account

	// Input sanitation
	if len(args) != 3 {
		return shim.Error("incorrect number of arguments. 3 expected")
	}
	if args[0] == "" {
		return shim.Error("1st argument must be a non-empty string")
	}
	if args[1] == "" {
		return shim.Error("2nd argument must be a non-empty string")
	}
	if args[2] == "" {
		return shim.Error("3rd argument must be a non-empty string")
	}
	if args[0] == args[1] {
		return shim.Error("the transfer must be between different accounts")
	}

	// Mapping args to variables
	payerAccNumber, err := strconv.Atoi(args[0])
	if err != nil {
		return shim.Error("1st argument must be a numeric string")
	}

	receiverAccNumber, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("2nd argument must be a numeric string")
	}

	transferValue, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("3rd argument must be a numeric string")
	}

	// Get payer account
	chaincodeName := "cc-account"
	chaincodeArgs := util.ToChaincodeArgs("GetByNumber", strconv.Itoa(payerAccNumber))

	// If `channel` is empty, the caller's channel is assumed.
	response := stub.InvokeChaincode(chaincodeName, chaincodeArgs, "")
	if response.Status != shim.OK {
		return shim.Error("failed to invoke `" + chaincodeName + "` chaincode: " + response.Message)
	}
	payerAccAsBytes := response.Payload[:]

	// Get receiver Account
	chaincodeArgs = util.ToChaincodeArgs("GetByNumber", strconv.Itoa(receiverAccNumber))
	response = stub.InvokeChaincode(chaincodeName, chaincodeArgs, "")
	if response.Status != shim.OK {
		return shim.Error("failed to invoke `" + chaincodeName + "` chaincode: " + response.Message)
	}
	receiverAccAsBytes := response.Payload[:]

	// Parse payer Account to Account object
	err = json.Unmarshal(payerAccAsBytes, &payerAcc)
	if err != nil {
		return shim.Error("cannot unmarshal payer account to JSON: " + err.Error())
	}

	// Check payer Account accBalance
	if payerAcc.AccountBalance < transferValue {
		fmt.Println("[DEBUG] insufficient funds. payerAcc.AccountBalance = " + strconv.Itoa(payerAcc.AccountBalance))
		return shim.Error("payer insufficient funds")
	}

	// Parse receiver Account to Account object
	err = json.Unmarshal(receiverAccAsBytes, &receiverAcc)
	if err != nil {
		return shim.Error("cannot unmarshal receiver account to JSON: " + err.Error())
	}

	// Transfer money
	payerAcc.AccountBalance = payerAcc.AccountBalance - transferValue
	receiverAcc.AccountBalance = receiverAcc.AccountBalance + transferValue

	// Parse Accounts back to JSON
	newPayerAccAsBytes, err := json.Marshal(payerAcc)
	if err != nil {
		return shim.Error("failed to marshal payer account object: " + err.Error())
	}
	newReceiverAccAsBytes, err := json.Marshal(receiverAcc)
	if err != nil {
		return shim.Error("failed to marshal receiver account object: " + err.Error())
	}

	// Update payer Account
	chaincodeArgs = util.ToChaincodeArgs("Update", string(newPayerAccAsBytes[:]))
	response = stub.InvokeChaincode(chaincodeName, chaincodeArgs, "")
	if response.Status != shim.OK {
		return shim.Error("could not update payer account: " + response.Message)
	}

	// Update receiver Account
	chaincodeArgs = util.ToChaincodeArgs("Update", string(newReceiverAccAsBytes[:]))
	response = stub.InvokeChaincode(chaincodeName, chaincodeArgs, "")
	if response.Status != shim.OK {
		return shim.Error("could not receiver payer account: " + response.Message)
	}

	fmt.Println("[DEBUG] end transfer.Money")
	return shim.Success(nil)
}
