/*
Package card provides services in the context of card asset.
*/
package card

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/common/util"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// Card structure with 3 properties. Structure tags are used by encoding/json library
type Card struct {
	ObjectType    string `json:"docType"`
	CardNumber    int    `json:"cardNumber"`
	AccountNumber string `json:"accountNumber"`
}

// Create - creates new card and stores into chaincode state
// params: cardNumber, AccountNumber
func Create(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fmt.Println("-- Starting card Create")

	// Input sanitation
	if len(args) != 2 {
		return shim.Error("Error: Incorrect number of arguments. 2 are expected!")
	}
	if len(args[0]) <= 0 {
		return shim.Error("Error: 1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("Error: 2nd argument must be a non-empty string")
	}

	// Mapping args to variables
	cardNumber, err := strconv.Atoi(args[0])
	if err != nil {
		return shim.Error("Error: 1st argument must be a numeric string")
	}
	accountNumber, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Error: 2nd argument must be a numeric string")
	}
	cardNumberStr := strconv.Itoa(cardNumber)

	// Check if it already exists
	cardAsBytes, err := stub.GetState("CARD" + cardNumberStr)
	if err != nil {
		return shim.Error("Error: Failed to get card data: " + err.Error())
	} else if cardAsBytes != nil {
		return shim.Error("Error: This card already exists: " + cardNumberStr)
	}

	// Check if account exists
	chaincodeName := "cc-account"
	queryArgs := util.ToChaincodeArgs("GetByNumber", strconv.Itoa(accountNumber))
	channelName := ""
	response := stub.InvokeChaincode(chaincodeName, queryArgs, channelName)
	if response.Status != shim.OK {
		return shim.Error("Error: Check if the chaincode name or the function name and parameters or account number are valid in InvokeChaincode!")
	}

	// Create card object and marshal to JSON
	objectType := "Card"
	card := &Card{objectType, cardNumber, strconv.Itoa(accountNumber)}
	cardJSONasBytes, err := json.Marshal(card)
	if response.Status != shim.OK {
		msg := string(response.Payload[:])
		return shim.Error(msg)
	}

	// Save card to state
	err = stub.PutState("CARD"+cardNumberStr, cardJSONasBytes)
	if err != nil {
		return shim.Error("Error: Could not put state of card: " + err.Error())
	}

	// Card saved and indexed. Return success
	fmt.Println("-- Ending card Create")
	return shim.Success([]byte("Card created!"))
}

// GetByNumber - Performs a query based on card number
// param: CardNumber
func GetByNumber(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fmt.Println("-- Starting card GetByNumber")
	var err error

	// Input sanitation
	if len(args) != 1 {
		return shim.Error("Error: Incorrect number of arguments. 1 are expected!")
	}

	// Mapping arg to variable
	cardNumber := args[0]

	// Get card state and check if it exists
	cardAsJSON, err := stub.GetState("CARD" + cardNumber)
	if err != nil {
		return shim.Error("Error: Failed to get state of account: " + cardNumber)
	} else if cardAsJSON == nil {
		return shim.Error("Error: Card " + cardNumber + " does not exist!")
	}

	fmt.Println("-- Ending GetCardByNumber")
	return shim.Success(cardAsJSON)
}

// GetAll - Get all cards in World State
// params: none
func GetAll(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Println("-- Starting card: GetAll")

	cardsIterator, err := stub.GetStateByRange("", "")
	defer cardsIterator.Close()
	if err != nil {
		return shim.Error("Error while querying ledger. Error: " + err.Error())
	}

	var records []string
	if cardsIterator.HasNext() {
		for cardsIterator.HasNext() {
			recordAsBytes, err := cardsIterator.Next()
			if err != nil {
				return shim.Error("Error while iterating through ledger. Error: " + err.Error())
			}

			records = append(records, string(recordAsBytes.Value[:]))
		}
	}

	fmt.Println(records)
	return shim.Success([]byte("Success"))
}
