package account

import (
	"encoding/json"
	"fmt"
	"strconv"

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

/*
 * ============================================================
 * InitAccounts - creates five Accounts and stores into chaincode state
 * params:
 * ============================================================
 */
func InitAccounts(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Println("-- Starting InitAccounts")

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
			return shim.Error("Error inserting accounts: " + err.Error())
		}

		fmt.Println("Added", accounts[i])
		i = i + 1
	}

	fmt.Println("-- Ending InitAccounts")
	return shim.Success([]byte("Accounts created!"))
}

/*
 * ============================================================
 * CreateAccount - creates new Account and stores into chaincode state
 * params: Account idAccount, accBalance, accOwner
 * ============================================================
 */
func CreateAccount(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fmt.Println("-- Starting CreateAccount")

	var err error

	// Input sanitation
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. 3 are expected!")
	}
	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return shim.Error("3rd argument must be a non-empty string")
	}

	// Mapping args to variables
	accNumber, err := strconv.Atoi(args[0])
	if err != nil {
		return shim.Error("1st argument must be a numeric string")
	}
	accNumberAsStr := args[0]

	accBalance, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("2nd argument must be a numeric string")
	}
	accOwner := args[2]

	// Get Account state and check if it already exists
	AccountAsBytes, err := stub.GetState("ACC" + accNumberAsStr)
	if err != nil {
		return shim.Error("Failed to get Account data: " + err.Error())
	} else if AccountAsBytes != nil {
		return shim.Error("This Account already exists: " + accNumberAsStr)
	}

	// Create Account object and marshal to JSON
	objectType := "Account"
	account := &Account{objectType, accNumber, accBalance, accOwner}
	accountJSONasBytes, _ := json.Marshal(account)

	// Save Account to state
	err = stub.PutState("ACC"+strconv.Itoa(accNumber), accountJSONasBytes)
	if err != nil {
		return shim.Error("Error inserting account: " + err.Error())
	}

	// Account saved and indexed. Return success
	fmt.Println("-- Ending CreateAccount")
	return shim.Success([]byte("Account created!"))
}

/*
 * ============================================================
 * GetAccountByNumber - Performs a query based on Account number
 * param: AccountNumber
 * ============================================================
 */
func GetAccountByNumber(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fmt.Println("-- Starting GetAccountByNumber")
	var err error

	// Input sanitation
	if len(args) != 1 {
		return shim.Error("GetAccountByNumber: Incorrect number of arguments. 1 are expected!")
	}

	// Mapping arg to variable
	accNumber := args[0]

	// Get Account state and check if it exists

	accountAsBytes, err := stub.GetState("ACC" + accNumber)
	fmt.Println(accountAsBytes)
	if err != nil {
		return shim.Error("GetAccountByNumber: Fail to get state of account: " + accNumber)
	} else if accountAsBytes == nil {
		return shim.Error("GetAccountByNumber: Account " + accNumber + " does not exist!")
	}

	fmt.Println("-- Ending GetAccountByNumber")
	return shim.Success(accountAsBytes)
}

/*
 * ============================================================
 * UpdateAccountByNumber - Updates (rewrites) an account
 * param: AccountID, Account as bytes
 * ============================================================
 */
func UpdateAccountByNumber(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fmt.Println("-- Starting UpdateAccountByNumber")
	var err error

	// Input sanitation
	if args[0] == "" {
		return shim.Error("UpdateAccountByNumber error: 1st argument must be a non-empty string")
	}
	if args[1] == "" {
		return shim.Error("UpdateAccountByNumber error: 2nd argument must be a non-empty string")
	}
	_, err = strconv.Atoi(args[0])
	if err != nil {
		return shim.Error("UpdateAccountByNumber error: 1st argument must be a numeric string")
	}

	// Mapping arg to variable
	accNumber := args[0]
	accAsString := args[1]

	// Update (rewrite) Account
	err = stub.PutState("ACC"+accNumber, []byte(accAsString))
	if err != nil {
		return shim.Error("Error: " + err.Error() + ", updating account: " + accNumber)
	}

	fmt.Println("-- Ending UpdateAccountByNumber")
	return shim.Success([]byte("Account updated!"))
}
