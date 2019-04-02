/*
Package query provides a set of auxiliary functions that are able to
get data from couchdb database.
*/
package query

import (
	"bytes"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// ConstructQueryResponseFromIterator - Constructs a JSON array containing query results from
// a given result iterator
func ConstructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) ([]byte, error) {
	fmt.Println("[DEBUG] begin query.ConstructQueryResponseFromIterator")

	// Buffer is a JSON array containing QueryResults
	var b bytes.Buffer
	b.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			b.WriteString(",")
		}
		b.WriteString("{\"Key\":")
		b.WriteString("\"")
		b.WriteString(queryResponse.Key)
		b.WriteString("\"")
		b.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		b.WriteString(string(queryResponse.Value[:]))
		b.WriteString("}")

		bArrayMemberAlreadyWritten = true
	}

	b.WriteString("]")

	fmt.Println("[DEBUG] end query.ConstructQueryResponseFromIterator")
	return b.Bytes(), nil
}

// GetQueryResultForQueryString - Executes the passed in query string.
// Result set is built and returned as a byte array containing the JSON results.
func GetQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {
	fmt.Println("[DEBUG] begin query.GetQueryResultForQueryString")
	fmt.Printf("[DEBUG] queryString:\n%s\n", queryString)

	var queryResult []byte

	// Query couchdb
	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}

	defer resultsIterator.Close()

	// Check if there are query results
	if resultsIterator.HasNext() {
		// Format query response
		queryResult, err = ConstructQueryResponseFromIterator(resultsIterator)
		if err != nil {
			return nil, err
		}
	}

	fmt.Printf("[DEBUG] queryResult:\n%s\n", queryResult)
	fmt.Println("[DEBUG] end query.GetQueryResultForQueryString")
	return queryResult, nil
}
