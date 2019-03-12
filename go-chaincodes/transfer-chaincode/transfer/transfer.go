package transfer

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/go-chaincodes/account-chaincode/account"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// toChaincodeArgs - Prepares function arguments to invoke
// params: args
func toChaincodeArgs(args ...string) [][]byte {
	bargs := make([][]byte, len(args))
	for i, arg := range args {
		bargs[i] = []byte(arg)
	}
	return bargs
}

// Money - Transfer money between Accounts
// param: AccountNumber, AccountNumber, Value
func Money(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fmt.Println("-- Starting TransferMoney")

	var err error

	// Input sanitation
	if len(args) != 3 {
		return shim.Error("Error: Incorrect number of arguments. 3 are expected!")
	}
	if len(args[0]) <= 0 {
		return shim.Error("Error: 1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("Error: 2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return shim.Error("Error: 3rd argument must be a non-empty string")
	}
	if args[0] == args[1] {
		return shim.Error("Error: The transfer must be between different accounts")
	}

	// Mapping args to variables
	payerAccNumber, err := strconv.Atoi(args[0])
	if err != nil {
		return shim.Error("Error: 1st argument must be a numeric string")
	}

	receiverAccNumber, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Error: 2nd argument must be a numeric string")
	}

	value, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("Error: 3rd argument must be a numeric string")
	}

	// Get payer account
	chaincodeName := "cc-account"
	queryArgs := toChaincodeArgs("GetByNumber", strconv.Itoa(payerAccNumber))
	channelName := ""
	response := stub.InvokeChaincode(chaincodeName, queryArgs, channelName)
	if response.Status != shim.OK {
		msg := string(response.Payload[:])
		return shim.Error(msg)
	}
	payerAccAsString := string(response.Payload[:])

	// Get receiver Account
	queryArgs = toChaincodeArgs("GetByNumber", strconv.Itoa(receiverAccNumber))
	response = stub.InvokeChaincode(chaincodeName, queryArgs, channelName)
	if response.Status != shim.OK {
		return shim.Error("Error: Check if the account number is valid!")
	}
	receiverAccAsString := string(response.Payload[:])

	// Parse payer Account to Account object
	payerAcc := account.Account{}
	err = json.Unmarshal([]byte(payerAccAsString), &payerAcc)
	if err != nil {
		return shim.Error("Error: Failed to parse payer Account JSON: " + err.Error())
	}

	// Check payer Account accBalance
	if payerAcc.AccountBalance < value {
		return shim.Error("Error: Payer has insuffient money. Current account balance: " + strconv.Itoa(payerAcc.AccountBalance))
	}

	// Parse receiver Account to Account object
	receiverAcc := account.Account{}
	err = json.Unmarshal([]byte(receiverAccAsString), &receiverAcc)
	if err != nil {
		return shim.Error("Error: Failed to parse receiver Account JSON: " + err.Error())
	}

	// Transfer money
	payerAcc.AccountBalance = payerAcc.AccountBalance - value
	receiverAcc.AccountBalance = receiverAcc.AccountBalance + value

	// Parse Accounts back to JSON
	payerAccAsBytes, err := json.Marshal(payerAcc)
	if err != nil {
		return shim.Error("Error: Failed to parse payer Account object: " + err.Error())
	}
	receiverAccAsBytes, err := json.Marshal(receiverAcc)
	if err != nil {
		return shim.Error("Error: Failed to parse receiver Account object: " + err.Error())
	}

	// Update payer Account
	payerAccAsString = string(payerAccAsBytes[:])
	queryArgs = toChaincodeArgs("UpdateByNumber", strconv.Itoa(payerAccNumber), payerAccAsString)
	response = stub.InvokeChaincode(chaincodeName, queryArgs, channelName)
	if response.Status != shim.OK {
		return shim.Error("Error: Could not update payer account!")
	}

	// Update receiver Account
	receiverAccAsString = string(receiverAccAsBytes[:])
	queryArgs = toChaincodeArgs("UpdateByNumber", strconv.Itoa(receiverAccNumber), receiverAccAsString)
	response = stub.InvokeChaincode(chaincodeName, queryArgs, channelName)
	if response.Status != shim.OK {
		return shim.Error("Error: Could not receiver payer account!")
	}

	fmt.Println("-- Ending TransferMoney")
	return shim.Success([]byte("Money transfered!"))
}
