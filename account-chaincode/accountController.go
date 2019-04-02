/*
==== Install/Instantiate/Upgrade
peer chaincode install -n cc-account -p github.com/hyperledger-fabric-go-chaincodes/account-chaincode -v v1
peer chaincode instantiate -o orderer.example.com:7050 -C mychannel -n cc-account -c '{"Args":["debug"]}' -v v1
peer chaincode upgrade -o orderer.example.com:7050 -C mychannel -n cc-account -c '{"Args":["debug"]}' -v v2

==== List chaincodes ====
peer chaincode list --installed
peer chaincode list --instantiated -C mychannel

==== Accounts ====
 +++ Invokes
peer chaincode invoke -C mychannel -n cc-account -c '{"Args":["Init"]}'
peer chaincode invoke -C mychannel -n cc-account -c '{"Args":["Create","1","1000","Elcius"]}'
peer chaincode invoke -C mychannel -n cc-account -c '{"Args":["Create","2","1000","Natan"]}'
peer chaincode invoke -C mychannel -n cc-account -c '{"Args":["Create","6","1000","Marcelo"]}'
peer chaincode invoke -C mychannel -n cc-account -c '{"Args":["Delete","1"]}'
peer chaincode invoke -C mychannel -n cc-account -c '{"Args":["Update","{\"accountBalance\":7000,\"accountNumber\":2,\"accountOwner\":\"Natanael\",\"docType\":\"Account\"}"]}'

 +++ Queries
peer chaincode query -C mychannel -n cc-account -c '{"Args":["GetAll"]}' | jq
peer chaincode query -C mychannel -n cc-account -c '{"Args":["GetByNumber","1"]}' | jq
peer chaincode query -C mychannel -n cc-account -c '{"Args":["GetByOwner","Elcius"]}' | jq
peer chaincode query -C mychannel -n cc-account -c '{"Args":["GetHistory","1"]}' | jq
*/

package main

import (
	"strings"

	"github.com/hyperledger-fabric-go-chaincodes/account-chaincode/account"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// AccountsChaincode struct
type AccountsChaincode struct {
}

// Logger
var logger = shim.NewLogger("cc-account")

// Main
func main() {
	err := shim.Start(new(AccountsChaincode))
	if err != nil {
		logger.SetLevel(shim.LogCritical)
		logger.Critical("Failed to initialize accounts chaincode: " + err.Error())
	}
}

// Init - initializes chaincode
func (t *AccountsChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	args := stub.GetStringArgs()
	var logLevel string

	// Input sanitation
	if len(args) > 1 {
		return shim.Error("Incorrect number of arguments. None or 1 expected")
	}

	// Input Mapping
	if len(args) == 1 {
		logLevel = strings.ToUpper(args[0])
	}

	// Selecting log level
	switch logLevel {
	case "DEBUG":
		logger.SetLevel(shim.LogDebug)
	case "INFO":
		logger.SetLevel(shim.LogInfo)
	case "NOTICE":
		logger.SetLevel(shim.LogNotice)
	case "WARNING":
		logger.SetLevel(shim.LogWarning)
	case "ERROR":
		logger.SetLevel(shim.LogError)
	case "CRITICAL":
		logger.SetLevel(shim.LogCritical)
	default:
		logger.SetLevel(shim.LogInfo)
		logger.Warning("Level \"" + logLevel + "\" not recognized as valid log level")
		logger.Notice("Using default logger level \"INFO\"")
	}

	logger.Info("Initialized `cc-account` chaincode")
	return shim.Success(nil)
}

// Invoke - Entry point for Invocations
func (t *AccountsChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()

	// Configuring logger
	// logger.SetLevel(shim.LogDebug)
	logger.Info("Chaincode invoke: function:\"" + function + "\"")

	// Handle different functions
	switch function {
	case "Init":
		return account.Init(stub, logger)
	case "Create":
		return account.Create(stub, logger, args)
	case "GetAll":
		return account.GetAll(stub, logger)
	case "GetByNumber":
		return account.GetByNumber(stub, logger, args)
	case "GetByOwner":
		return account.GetByOwner(stub, logger, args)
	case "Update":
		return account.Update(stub, logger, args)
	case "Delete":
		return account.Delete(stub, logger, args)
	case "GetHistory":
		return account.GetHistoryByAccNumber(stub, logger, args)
	default:
		// Error
		logger.Error("Received unknown function invoke: \"" + function + "\"")
		return shim.Error("Received unknown function invoke: \"" + function + "\"")
	}
}
