/*
Package account provides services in the context of account asset.
*/
package account

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger-fabric-go-chaincodes/query"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// Account structure with 4 properties. Structure tags are used by encoding/json library
type Account struct {
	ObjectType     string `json:"docType"`
	AccountNumber  int    `json:"accountNumber"`
	AccountBalance int    `json:"accountBalance"`
	AccountOwner   string `json:"accountOwner"`
}

// Init - creates five Accounts and stores into chaincode state
// params: none
func Init(stub shim.ChaincodeStubInterface, logger *shim.ChaincodeLogger) peer.Response {
	logger.Info("Entry method: Init")

	accounts := []Account{
		{ObjectType: "Account", AccountNumber: 1, AccountBalance: 1000, AccountOwner: "Elcius"},
		{ObjectType: "Account", AccountNumber: 2, AccountBalance: 1000, AccountOwner: "Natan"},
		{ObjectType: "Account", AccountNumber: 3, AccountBalance: 1000, AccountOwner: "Johan"},
		{ObjectType: "Account", AccountNumber: 4, AccountBalance: 1000, AccountOwner: "Leandro"},
		{ObjectType: "Account", AccountNumber: 5, AccountBalance: 1000, AccountOwner: "Marcos"},
	}

	for i := 0; i < len(accounts); i++ {
		accountsAsBytes, _ := json.Marshal(accounts[i])

		err := stub.PutState("ACC"+strconv.Itoa(i+1), accountsAsBytes)
		if err != nil {
			logger.Error("Error inserting accounts:", err.Error())
			logger.Info("Exit method: Init")
			return shim.Error("Error inserting accounts: " + err.Error())
		}

		logger.Debug("pushed ACC" + strconv.Itoa(i+1) + ":", accounts[i])
	}

	err := stub.SetEvent("accounts_created", []byte("Success"))
	if err != nil {
		logger.Critical("Failed to set event `accounts_created`:", err.Error())
		logger.Info("Exit method: Init")
		return shim.Error("Failed to set event `accounts_created`: " +  err.Error())
	}

	logger.Info("Exit method: Init")
	return shim.Success(nil)
}

// Create - creates new Account and stores into chaincode state
// params: Account idAccount, accBalance, accOwner
func Create(stub shim.ChaincodeStubInterface, logger *shim.ChaincodeLogger, args []string) peer.Response {
	logger.Info("Entry method: Create")
	logger.Debug("Received args:", args)

	var err error

	// Input sanitation
	if len(args) != 3 {
		logger.Info("Exit method: Create")
		return shim.Error("incorrect number of arguments. 3 expected")
	}
	if args[0] == "" {
		logger.Info("Exit method: Create")
		return shim.Error("1st argument must be a non-empty string")
	}
	if args[1] == "" {
		logger.Info("Exit method: Create")
		return shim.Error("2nd argument must be a non-empty string")
	}
	if args[2] == "" {
		logger.Info("Exit method: Create")
		return shim.Error("3rd argument must be a non-empty string")
	}

	// Mapping args to variables
	accNumberAsStr := args[0]
	accNumber, err := strconv.Atoi(accNumberAsStr)
	if err != nil {
		logger.Info("Exit method: Create")
		return shim.Error("1st argument must be a numeric string")
	}

	accBalance, err := strconv.Atoi(args[1])
	if err != nil {
		logger.Info("Exit method: Create")
		return shim.Error("2nd argument must be a numeric string")
	}

	accOwner := args[2]

	// Get Account state and check if it already exists
	AccountAsBytes, err := stub.GetState("ACC" + accNumberAsStr)
	if err != nil {
		logger.Info("Exit method: Create")
		return shim.Error("Failed to get account data: " + err.Error())
	} else if AccountAsBytes != nil {
		logger.Info("Exit method: Create")
		return shim.Error("Account ACC" + accNumberAsStr + " already exists")
	}

	// Create Account object and marshal to JSON
	objectType := "Account"
	account := &Account{objectType, accNumber, accBalance, accOwner}
	accountJSONasBytes, err := json.Marshal(account)
	if err != nil {
		logger.Info("Exit method: Create")
		return shim.Error("Cannot marshal Account: " + err.Error())
	}

	// Save Account to state
	err = stub.PutState("ACC"+accNumberAsStr, accountJSONasBytes)
	if err != nil {
		logger.Info("Exit method: Create")
		return shim.Error("Failed to put state of account: " + err.Error())
	}

	// Account saved and indexed. Return success
	err = stub.SetEvent("account_created", []byte("Success"))
	if err != nil {
		logger.Critical("Failed to set event `account_created`: " + err.Error())
		logger.Info("Exit method: Create")
		return shim.Error("Failed to set event `account_created`: " + err.Error())
	}

	logger.Info("Exit method: Create")
	return shim.Success(nil)
}

// GetAll - Get all the existing accounts
// params: none needed
func GetAll(stub shim.ChaincodeStubInterface, logger *shim.ChaincodeLogger) peer.Response {
	logger.Info("Entry method: GetAll")
	var err error

	// Get Account state and check if it exists
	accountsIterator, err := stub.GetStateByRange("", "")
	if err != nil {
		logger.Info("Exit method: GetAll")
		return shim.Error("Cannot get ledger state: " + err.Error())
	}

	defer accountsIterator.Close()

	queryResults, err := query.ConstructQueryResponseFromIterator(accountsIterator)
	if err != nil {
		logger.Info("Exit method: GetAll")
		return shim.Error("Failed to construct results from iterator: " + err.Error())
	}

	logger.Debug("queryResults: " + string(queryResults[:]))

	err = stub.SetEvent("get_all_accounts", []byte("Success"))
	if err != nil {
		logger.Critical("Failed to set event `get_all_accounts`: " + err.Error())
		logger.Info("Exit method: GetAll")
		return shim.Error("Failed to set event `get_all_accounts`: " + err.Error())
	}

	logger.Info("Exit method: GetAll")
	return shim.Success(queryResults)
}

// GetByNumber - Performs a query based on Account number
// param: AccountNumber
func GetByNumber(stub shim.ChaincodeStubInterface, logger *shim.ChaincodeLogger, args []string) peer.Response {
	logger.Info("Entry method: GetByNumber")
	logger.Debug("Received args:", args)

	// Input sanitation
	if len(args) != 1 {
		logger.Info("Exit method: GetByNumber")
		return shim.Error("Incorrect number of arguments. 1 expected")
	}
	if args[0] == "" {
		logger.Info("Exit method: GetByNumber")
		return shim.Error("Account number must be a non-empty string")
	}
	_, err := strconv.Atoi(args[0])
	if err != nil {
		logger.Info("Exit method: GetByNumber")
		return shim.Error("Account number must be numeric string")
	}

	// Mapping arg to variable
	accNumber := args[0]

	// Get Account state and check if it exists
	accountAsBytes, err := stub.GetState("ACC" + accNumber)
	if err != nil {
		logger.Info("Exit method: GetByNumber")
		return shim.Error("Failed to fetch account ACC" + accNumber + " from ledger: " + err.Error())
	} else if accountAsBytes == nil {
		logger.Info("Exit method: GetByNumber")
		return shim.Error("Account ACC" + accNumber + " does not exist")
	}

	err = stub.SetEvent("get_account_by_number", []byte("Success"))
	if err != nil {
		logger.Critical("Failed to set event `get_account_by_number`: " + err.Error())
		logger.Info("Exit method: GetByNumber")
		return shim.Error("Failed to set event `get_account_by_number`: " + err.Error())
	}

	logger.Info("Exit method: GetByNumber")
	return shim.Success(accountAsBytes)
}

// GetByOwner - Queries account by the owner name
// param: accountOwner
func GetByOwner(stub shim.ChaincodeStubInterface, logger *shim.ChaincodeLogger, args []string) peer.Response {
	logger.Info("Entry method: GetByOwner")
	logger.Debug("Received args:", args)

	// Input sanitation
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. 1 expected")
	}
	if args[0] == "" {
		return shim.Error("Argument must be a non-empty string")
	}

	// Mapping arg to variable
	accOwner := args[0]

	// Construct query string using account owner name
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"Account\",\"accountOwner\":\"%s\"}}", accOwner)
	logger.Debug("Query string:", queryString)

	// Use package query to query couchdb and format the result
	queryResults, err := query.GetQueryResultForQueryString(stub, queryString)
	if err != nil {
		logger.Info("Exit method: GetByOwner")
		return shim.Error("Cannot get query results: " + err.Error())
	}

	err = stub.SetEvent("get_account_by_owner", []byte("Success"))
	if err != nil {
		logger.Critical("Failed to set event `get_account_by_owner`: " + err.Error())
		logger.Info("Exit method: GetByOwner")
		return shim.Error("Failed to set event `get_account_by_owner`: " + err.Error())
	}

	logger.Info("Exit method: GetByOwner")
	return shim.Success(queryResults)
}

// Update - Updates (rewrites) an account
// param: Account JSON as bytes
func Update(stub shim.ChaincodeStubInterface, logger *shim.ChaincodeLogger, args []string) peer.Response {
	logger.Info("Entry method: Update")
	logger.Debug("Received args:", args)

	var err error

	// Input sanitation
	if len(args) != 1 {
		logger.Info("Exit method: Update")
		return shim.Error("Incorrect number of arguments. 1 expected")
	}
	if args[0] == "" {
		logger.Info("Exit method: Update")
		return shim.Error("Argument must be a non-empty string")
	}

	// Mapping arg to variable
	accAsString := args[0]

	// Validating string input
	var accObject Account
	err = json.Unmarshal([]byte(accAsString), &accObject)
	if err != nil {
		logger.Info("Exit method: Update")
		return shim.Error("Account not valid as json object: " + err.Error())
	}

	// Update (rewrite) Account
	accNumber := strconv.Itoa(accObject.AccountNumber)
	err = stub.PutState("ACC"+accNumber, []byte(accAsString))
	if err != nil {
		logger.Info("Exit method: Update")
		return shim.Error("Failed to update ACC" + accNumber + ": " + err.Error())
	}

	err = stub.SetEvent("update_account", []byte("Success"))
	if err != nil {
		logger.Critical("Failed to set event `update_account`: " + err.Error())
		logger.Info("Exit method: Update")
		return shim.Error("Failed to set event `update_account`: " + err.Error())
	}

	logger.Info("Exit method: Update")
	return shim.Success(nil)
}

// Delete - Delete account based on its number
// param: AccountNumber
func Delete(stub shim.ChaincodeStubInterface, logger *shim.ChaincodeLogger, args []string) peer.Response {
	logger.Info("Entry method: Delete")
	logger.Debug("Received args:", args)

	var err error

	// Input sanitation
	if args[0] == "" {
		logger.Info("Exit method: Delete")
		return shim.Error("1st argument must be a non-empty string")
	}
	_, err = strconv.Atoi(args[0])
	if err != nil {
		logger.Info("Exit method: Delete")
		return shim.Error("1st argument must be a numeric string")
	}

	// Mapping arg to variable
	accNumber := args[0]

	// Get Account state and check if it exists
	accountAsBytes, err := stub.GetState("ACC" + accNumber)
	if err != nil {
		logger.Info("Exit method: Delete")
		return shim.Error("Failed to fetch account ACC" + accNumber + " from ledger: " + err.Error())
	} else if accountAsBytes == nil {
		logger.Info("Exit method: Delete")
		return shim.Error("Account ACC" + accNumber + " does not exist")
	}

	// Remove the account from chaincode state
	err = stub.DelState("ACC" + accNumber)
	if err != nil {
		logger.Info("Exit method: Delete")
		return shim.Error("Failed to delete state: " + err.Error())
	}

	err = stub.SetEvent("delete_account", []byte("Success"))
	if err != nil {
		logger.Critical("Failed to set event `delete_account`: " + err.Error())
		logger.Info("Exit method: Delete")
		return shim.Error("Failed to set event `delete_account`: " + err.Error())
	}

	logger.Info("Exit method: Delete")
	return shim.Success(nil)
}

// GetHistory - Queries the history for a given account and returns on JSON format
// param: AccountNumber
func GetHistoryByAccNumber(stub shim.ChaincodeStubInterface, logger *shim.ChaincodeLogger, args []string) peer.Response {
	logger.Info("Entry method: GetHistory")
	logger.Debug("Received args:", args)

	var err error
	var b bytes.Buffer

	// Input sanitation
	if len(args) != 1 {
		logger.Info("Exit method: GetHistory")
		return shim.Error("Incorrect number of arguments. 1 expected")
	}
	_, err = strconv.Atoi(args[0])
	if err != nil {
		logger.Info("Exit method: GetHistory")
		return shim.Error("Argument must be a numeric string")
	}

	// Mapping arg to variable
	accNumber := args[0]

	// Get History iterator
	resultsIterator, err := stub.GetHistoryForKey("ACC" + accNumber)
	if err != nil {
		logger.Info("Exit method: GetHistory")
		return shim.Error("Failed to fetch asset history: " + err.Error())
	}
	defer resultsIterator.Close()

	bArrayMemberAlreadyWritten := false

	if resultsIterator.HasNext() {
		b.WriteString("[")

		// Itarate over results
		for resultsIterator.HasNext() {
			historyData, err := resultsIterator.Next()
			if err != nil {
				return shim.Error("failed to iterate over results: " + err.Error())
			}

			// Add a comma before array members, suppress it for the first array member
			if bArrayMemberAlreadyWritten == true {
				b.WriteString(",")
			}

			b.WriteString("{\"TxID\":")
			b.WriteString("\"")
			b.WriteString(historyData.TxId)
			b.WriteString("\"")
			b.WriteString(", \"Value\":")

			// Check if account has been deleted
			if historyData.Value == nil {
				b.WriteString("{}")
				b.WriteString(", \"IsDeleted\":")
				b.WriteString("true")
				b.WriteString("}")
			} else {
				b.Write(historyData.Value)
				b.WriteString(", \"IsDeleted\":")
				b.WriteString("false")
				b.WriteString("}")
			}
			bArrayMemberAlreadyWritten = true
		}

		b.WriteString("]")
	} else {
		logger.Info("Exit method: GetHistory")
		return shim.Error("Cannot find account history. ACC" + accNumber + " does not exist")
	}

	err = stub.SetEvent("get_history", []byte("Success"))
	if err != nil {
		logger.Critical("Failed to set event `get_history`: " + err.Error())
		logger.Info("Exit method: GetHistory")
		return shim.Error("Failed to set event `get_history`: " + err.Error())
	}

	logger.Info("Exit method: GetHistory")
	return shim.Success(b.Bytes())
}
