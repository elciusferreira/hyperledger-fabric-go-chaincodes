package account

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/go-chaincodes/utils/qryutils"
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
	fmt.Println("-- Starting account Init")

	txID := stub.GetTxID()
	fmt.Println("Transaction ID:", txID)

	accounts := []Account{
		Account{ObjectType: "Account", AccountNumber: 1, AccountBalance: 1000, AccountOwner: "Elcius"},
		Account{ObjectType: "Account", AccountNumber: 2, AccountBalance: 1000, AccountOwner: "Natan"},
		Account{ObjectType: "Account", AccountNumber: 3, AccountBalance: 1000, AccountOwner: "Johan"},
		Account{ObjectType: "Account", AccountNumber: 4, AccountBalance: 1000, AccountOwner: "Leandro"},
		Account{ObjectType: "Account", AccountNumber: 5, AccountBalance: 1000, AccountOwner: "Marcos"},
	}

	i := 0
	var err error
	for i < len(accounts) {
		fmt.Println("i is ", i)

		accountsAsBytes, _ := json.Marshal(accounts[i])
		fmt.Println("ACC" + strconv.Itoa(i+1))
		err = stub.PutState("ACC"+strconv.Itoa(i+1), accountsAsBytes)
		if err != nil {
			return shim.Error("Error: Failed to put state of accounts: " + err.Error())
		}

		fmt.Println("Added", accounts[i])
		i = i + 1
	}

	fmt.Println("-- Ending account Init")
	stub.SetEvent("accounts_created", []byte("Success"))
	return shim.Success([]byte("Accounts created!"))
}

// Create - creates new Account and stores into chaincode state
// params: Account idAccount, accBalance, accOwner
func Create(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fmt.Println("-- Starting account Create")

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

	// Mapping args to variables
	accNumber, err := strconv.Atoi(args[0])
	if err != nil {
		return shim.Error("Error: 1st argument must be a numeric string")
	}
	accNumberAsStr := args[0]

	accBalance, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Error: 2nd argument must be a numeric string")
	}
	accOwner := args[2]

	// Get Account state and check if it already exists
	AccountAsBytes, err := stub.GetState("ACC" + accNumberAsStr)
	if err != nil {
		return shim.Error("Error: Failed to get Account data: " + err.Error())
	} else if AccountAsBytes != nil {
		return shim.Error("Error: This Account already exists: " + accNumberAsStr)
	}

	// Create Account object and marshal to JSON
	objectType := "Account"
	account := &Account{objectType, accNumber, accBalance, accOwner}
	accountJSONasBytes, _ := json.Marshal(account)

	// Save Account to state
	err = stub.PutState("ACC"+strconv.Itoa(accNumber), accountJSONasBytes)
	if err != nil {
		return shim.Error("Error: Failed to put state of account: " + err.Error())
	}

	// Account saved and indexed. Return success
	fmt.Println("-- Ending account Create")
	return shim.Success([]byte("Account created!"))
}

// GetByNumber - Performs a query based on Account number
// param: AccountNumber
func GetByNumber(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fmt.Println("-- Starting account GetByNumber")
	var err error

	// Input sanitation
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. 1 are expected!")
	}

	// Mapping arg to variable
	accNumber := args[0]

	// Get Account state and check if it exists
	accountAsBytes, err := stub.GetState("ACC" + accNumber)
	if err != nil {
		return shim.Error("Error: Fail to get state of account: " + accNumber)
	} else if accountAsBytes == nil {
		return shim.Error("Error: Account " + accNumber + " does not exist!")
	}

	fmt.Println("-- Ending account GetByNumber")
	return shim.Success(accountAsBytes)
}

// GetByOwner - Queries account by the owner name
// param: accountOwner
func GetByOwner(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fmt.Println("-- Starting account GetByOwner")

	// Input sanitation
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	if args[0] == "" {
		return shim.Error("Error: Argument must be a non-empty string")
	}

	// Mapping arg to variable
	accOwner := args[0]

	// Construct query string using account owner name
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"Account\",\"accountOwner\":\"%s\"}}", accOwner)

	// Use package qryutils to query couchdb and format the result
	queryResults, status := qryutils.GetQueryResultForQueryString(stub, queryString)
	if status != "SUCCESS" {
		return shim.Error(status)
	}

	fmt.Println("-- Ending account GetByOwner")
	return shim.Success([]byte(queryResults))
}

// UpdateByNumber - Updates (rewrites) an account
// param: AccountID, Account as bytes
func UpdateByNumber(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fmt.Println("-- Starting account UpdateByNumber")
	var err error

	// Input sanitation
	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	if args[0] == "" {
		return shim.Error("Error: 1st argument must be a non-empty string")
	}
	if args[1] == "" {
		return shim.Error("Error: 2nd argument must be a non-empty string")
	}
	_, err = strconv.Atoi(args[0])
	if err != nil {
		return shim.Error("Error: 1st argument must be a numeric string")
	}

	// Mapping arg to variable
	accNumber := args[0]
	accAsString := args[1]

	// Update (rewrite) Account
	err = stub.PutState("ACC"+accNumber, []byte(accAsString))
	if err != nil {
		return shim.Error("Error: " + err.Error() + ", updating account: " + accNumber)
	}

	fmt.Println("-- Ending account UpdateByNumber")
	return shim.Success([]byte("Account updated!"))
}

// Delete - Delete account based on its number
// param: AccountNumber
func Delete(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fmt.Println("-- Starting account Delete")

	var err error

	// Input sanitation
	if args[0] == "" {
		return shim.Error("Error: 1st argument must be a non-empty string")
	}
	_, err = strconv.Atoi(args[0])
	if err != nil {
		return shim.Error("Error: 1st argument must be a numeric string")
	}

	// Mapping arg to variable
	accNumber := args[0]

	// Get Account state and check if it exists
	accountAsBytes, err := stub.GetState("ACC" + accNumber)
	if err != nil {
		return shim.Error("Error: Fail to get state of account: " + accNumber)
	} else if accountAsBytes == nil {
		return shim.Error("Error: Account " + accNumber + " does not exist!")
	}

	// Remove the account from chaincode state
	err = stub.DelState("ACC" + accNumber)
	if err != nil {
		return shim.Error("Error: Failed to delete state:" + err.Error())
	}

	fmt.Println("-- Ending account Delete")
	return shim.Success([]byte("Account deleted!"))
}

// GetHistory - Queries the history for a given account and returns on JSON format
// param: AccountNumber
func GetHistory(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fmt.Println("-- Starting GetHistory")

	var err error

	// Input sanitation
	if len(args) != 1 {
		return shim.Error("Error: Incorrect number of arguments. Expecting 1")
	}
	_, err = strconv.Atoi(args[0])
	if err != nil {
		return shim.Error("Error: 1st argument must be a numeric string")
	}

	// Mapping arg to variable
	accNumber := args[0]

	// Auxiliary struct
	type AuditHistory struct {
		TxID      string  `json:"TxID"`
		Value     Account `json:"Value"`
		IsDeleted bool    `json:"IsDeleted"`
	}

	var account Account
	// Store all transactions ID and account states
	var history []AuditHistory

	// Get History
	resultsIterator, err := stub.GetHistoryForKey("ACC" + accNumber)
	if err != nil {
		return shim.Error("Error: Failed to get history: " + err.Error())
	}
	defer resultsIterator.Close()

	if resultsIterator.HasNext() {
		// Itarate over results
		for resultsIterator.HasNext() {
			historyData, err := resultsIterator.Next()
			if err != nil {
				return shim.Error("Error: Failed to iterate over results: " + err.Error())
			}

			var tx AuditHistory

			// Copy transaction id over
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
		return shim.Error("Error: Failed to find history. Account " + accNumber + " does not exist!")
	}

	// Parse history list to array of bytes
	historyAsBytes, _ := json.Marshal(history)

	fmt.Printf("-- Ending GetHistory")
	return shim.Success(historyAsBytes)
}
