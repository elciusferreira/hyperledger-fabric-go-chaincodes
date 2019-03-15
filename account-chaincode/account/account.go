/*
Package account provides services in the context of account asset.
*/
package account

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/go-chaincodes/query"
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
func Init(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Println("[DEBUG] begin account.Init")

	txID := stub.GetTxID()
	fmt.Println("[DEBUG] Transaction ID:", txID)

	accounts := []Account{
		{ObjectType: "Account", AccountNumber: 1, AccountBalance: 1000, AccountOwner: "Elcius"},
		{ObjectType: "Account", AccountNumber: 2, AccountBalance: 1000, AccountOwner: "Natan"},
		{ObjectType: "Account", AccountNumber: 3, AccountBalance: 1000, AccountOwner: "Johan"},
		{ObjectType: "Account", AccountNumber: 4, AccountBalance: 1000, AccountOwner: "Leandro"},
		{ObjectType: "Account", AccountNumber: 5, AccountBalance: 1000, AccountOwner: "Marcos"},
	}

	for i := 0; i < len(accounts); i++ {
		accountsAsBytes, _ := json.Marshal(accounts[i])

		fmt.Println("[DEBUG] Key: ACC" + strconv.Itoa(i+1))

		err := stub.PutState("ACC" + strconv.Itoa(i+1), accountsAsBytes)
		if err != nil {
			return shim.Error("error inserting accounts: " + err.Error())
		}

		fmt.Println("[DEBUG] Added: ", accounts[i])
	}

	fmt.Println("[DEBUG] end account.Init")
    stub.SetEvent("accounts_created", []byte("Success"))
	return shim.Success(nil)
}

// Create - creates new Account and stores into chaincode state
// params: Account idAccount, accBalance, accOwner
func Create(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fmt.Println("[DEBUG] begin account.Create")

	var err error

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

	// Mapping args to variables
	accNumberAsStr := args[0]
	accNumber, err := strconv.Atoi(accNumberAsStr)
	if err != nil {
		return shim.Error("1st argument must be a numeric string")
	}

	accBalance, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("2nd argument must be a numeric string")
	}

	accOwner := args[2]

	// Get Account state and check if it already exists
	AccountAsBytes, err := stub.GetState("ACC" + accNumberAsStr)
	if err != nil {
		return shim.Error("failed to get account data: " + err.Error())
	} else if AccountAsBytes != nil {
		return shim.Error("account ACC" + accNumberAsStr + " already exists")
	}

	// Create Account object and marshal to JSON
	objectType := "Account"
	account := &Account{objectType, accNumber, accBalance, accOwner}
	accountJSONasBytes, err := json.Marshal(account)
	if err != nil {
		return shim.Error("cannot marshal Account: " + err.Error())
	}

	// Save Account to state
	err = stub.PutState("ACC" + accNumberAsStr, accountJSONasBytes)
	if err != nil {
		return shim.Error("failed to put state of account: " + err.Error())
	}

	// Account saved and indexed. Return success
	fmt.Println("[DEBUG] end account.Create")
	//return shim.Success([]byte("Account created!"))
	return shim.Success(nil)
}

// GetAll - Get all the existing accounts
// params: none needed
func GetAll(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Println("[DEBUG] begin account.GetAll")
	var err error

	// Get Account state and check if it exists
	accountsIterator, err := stub.GetStateByRange("", "")
	if err != nil {
		return shim.Error("cannot get ledger state: " + err.Error())
	}

	defer accountsIterator.Close()

	queryResults, err := query.ConstructQueryResponseFromIterator(accountsIterator)
	if err != nil {
		return shim.Error("failed to construct results from iterator: " + err.Error())
	}

	fmt.Println("[DEBUG] queryResults: " + string(queryResults[:]))
	fmt.Println("[DEBUG] end account.GetByNumber")
	return shim.Success(queryResults)
}

// GetByNumber - Performs a query based on Account number
// param: AccountNumber
func GetByNumber(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fmt.Println("[DEBUG] begin account.GetByNumber")

	// Input sanitation
	if len(args) != 1 {
		return shim.Error("incorrect number of arguments. 1 expected")
	}
	if args[0] == "" {
		return shim.Error("account number must be a non-empty string")
	}
	_, err := strconv.Atoi(args[0])
	if err != nil {
		return shim.Error("account number must be numeric string")
	}

	// Mapping arg to variable
	accNumber := args[0]

	// Get Account state and check if it exists
	accountAsBytes, err := stub.GetState("ACC" + accNumber)
	if err != nil {
		return shim.Error("failed to fetch account ACC" + accNumber + " from ledger: " + err.Error())
	} else if accountAsBytes == nil {
		return shim.Error("account ACC" + accNumber + " does not exist")
	}

	fmt.Println("[DEBUG] end account.GetByNumber")
	return shim.Success(accountAsBytes)
}

// GetByOwner - Queries account by the owner name
// param: accountOwner
func GetByOwner(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fmt.Println("[DEBUG] begin account.GetByOwner")

	// Input sanitation
	if len(args) != 1 {
		return shim.Error("incorrect number of arguments. 1 expected")
	}
	if args[0] == "" {
		return shim.Error("argument must be a non-empty string")
	}

	// Mapping arg to variable
	accOwner := args[0]

	// Construct query string using account owner name
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"Account\",\"accountOwner\":\"%s\"}}", accOwner)

	// Use package query to query couchdb and format the result
	queryResults, err := query.GetQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error("cannot get query results: " + err.Error())
	}

	fmt.Println("[DEBUG] end account.GetByOwner")
	return shim.Success(queryResults)
}

/* DEPRECATED
// UpdateByNumber - Updates (rewrites) an account
// param: AccountID, Account as bytes
func UpdateByNumber(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fmt.Println("[DEBUG] begin account.UpdateByNumber")
	var err error

	// Input sanitation
	if len(args) != 2 {
		return shim.Error("incorrect number of arguments. 2 expected")
	}
	if args[0] == "" {
		return shim.Error("1st argument must be a non-empty string")
	}
	if args[1] == "" {
		return shim.Error("2nd argument must be a non-empty string")
	}
	_, err = strconv.Atoi(args[0])
	if err != nil {
		return shim.Error("1st argument must be a numeric string")
	}

	// Mapping arg to variable
	accNumber := args[0]
	accAsString := args[1]

	// Validating account input
	var inputAcc Account
	err = json.Unmarshal([]byte(accAsString), &inputAcc)
	if err != nil {
		return shim.Error("account not valid as json: " + err.Error())
	}

	// Update (rewrite) Account
	err = stub.PutState("ACC" + accNumber, []byte(accAsString))
	if err != nil {
		return shim.Error("failed to update ACC" + accNumber + ": " + err.Error())
	}

	fmt.Println("[DEBUG] end account.UpdateByNumber")
	return shim.Success([]byte("Account updated"))
}
*/

// Update - Updates (rewrites) an account
// param: Account JSON as bytes
func Update(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fmt.Println("[DEBUG] begin account.Update")
	var err error

	// Input sanitation
	if len(args) != 1 {
		return shim.Error("incorrect number of arguments. 1 expected")
	}
	if args[0] == "" {
		return shim.Error("argument must be a non-empty string")
	}

	// Mapping arg to variable
	accAsString := args[0]

	// Validating string input
	var accObject Account
	err = json.Unmarshal([]byte(accAsString), &accObject)
	if err != nil {
		return shim.Error("account not valid as json object: " + err.Error())
	}

	// Update (rewrite) Account
	accNumber := strconv.Itoa(accObject.AccountNumber)
	err = stub.PutState("ACC" + accNumber, []byte(accAsString))
	if err != nil {
		return shim.Error("failed to update ACC" + accNumber + ": " + err.Error())
	}

	fmt.Println("[DEBUG] end account.Update")
	return shim.Success(nil)
}

// Delete - Delete account based on its number
// param: AccountNumber
func Delete(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fmt.Println("[DEBUG] begin account.Delete")

	var err error

	// Input sanitation
	if args[0] == "" {
		return shim.Error("1st argument must be a non-empty string")
	}
	_, err = strconv.Atoi(args[0])
	if err != nil {
		return shim.Error("1st argument must be a numeric string")
	}

	// Mapping arg to variable
	accNumber := args[0]

	// Get Account state and check if it exists
	accountAsBytes, err := stub.GetState("ACC" + accNumber)
	if err != nil {
		return shim.Error("failed to fetch account ACC" + accNumber + " from ledger: " + err.Error())
	} else if accountAsBytes == nil {
		return shim.Error("account ACC" + accNumber + " does not exist")
	}

	// Remove the account from chaincode state
	err = stub.DelState("ACC" + accNumber)
	if err != nil {
		return shim.Error("failed to delete state: " + err.Error())
	}

	fmt.Println("[DEBUG] end account.Delete")
	return shim.Success(nil)
}

// GetHistory - Queries the history for a given account and returns on JSON format
// param: AccountNumber
func GetHistoryByAccNumber(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fmt.Println("[DEBUG] begin account.GetHistory")

	var err error

	// Input sanitation
	if len(args) != 1 {
		return shim.Error("incorrect number of arguments. 1 expected")
	}
	_, err = strconv.Atoi(args[0])
	if err != nil {
		return shim.Error("argument must be a numeric string")
	}

	// Mapping arg to variable
	accNumber := args[0]

	// Auxiliary struct
	type AuditHistory struct {
		TxID      string  `json:"TxID"`
		Value     Account `json:"Value"`
		IsDeleted bool    `json:"IsDeleted"`
	}

	// Store all transactions ID and account states
	var history []AuditHistory
	var account Account

	// Get History
	resultsIterator, err := stub.GetHistoryForKey("ACC" + accNumber)
	if err != nil {
		return shim.Error("failed to fetch asset history: " + err.Error())
	}
	defer resultsIterator.Close()

	if resultsIterator.HasNext() {
		// Itarate over results
		for resultsIterator.HasNext() {
			historyData, err := resultsIterator.Next()
			if err != nil {
				return shim.Error("failed to iterate over results: " + err.Error())
			}

			// Copy transaction id over
			var tx AuditHistory
			tx.TxID = historyData.TxId

			// Check if account has been deleted
			if historyData.Value == nil {
				var emptyAccount Account
				// Copy nil account
				tx.Value = emptyAccount
				tx.IsDeleted = true
			} else {
				// Parse asset value to account object
				json.Unmarshal(historyData.Value, &account)

				// Copy account
				tx.Value = account
				tx.IsDeleted = false
			}
			// Add transaction (txID and account state or value) to the list
			history = append(history, tx)
		}
	} else {
		return shim.Error("cannot find account history. ACC" + accNumber + " does not exist")
	}

	// Parse history list to array of bytes
	historyAsBytes, _ := json.Marshal(history)

	fmt.Printf("[DEBUG] end account.GetHistory")
	return shim.Success(historyAsBytes)
}
