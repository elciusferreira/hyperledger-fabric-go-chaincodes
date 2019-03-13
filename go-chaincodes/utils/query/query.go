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

// constructQueryResponseFromIterator - Constructs a JSON array containing query results from
// a given result iterator
func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) (*bytes.Buffer, string) {
	fmt.Println("- Starting query constructQueryResponseFromIterator")

	// Buffer is a JSON array containing QueryResults
	var b bytes.Buffer
	b.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, "Error: Failed to get next result in results iterator"
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
		b.WriteString(string(queryResponse.Value))
		b.WriteString("}")

		bArrayMemberAlreadyWritten = true
	}

	b.WriteString("]")

	fmt.Println("- Ending query constructQueryResponseFromIterator")
	return &b, "SUCCESS"
}

// GetQueryResultForQueryString - Executes the passed in query string.
// Result set is built and returned as a byte array containing the JSON results.
func GetQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) (string, string) {
	fmt.Printf("- Starting query GetQueryResultForQueryString queryString:\n%s\n", queryString)

	var b *bytes.Buffer
	var status string

	// Query couchdb
	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return "ERROR", "Error: Failed to get query result: " + err.Error()
	}
	defer resultsIterator.Close()

	// Check if there are query results
	if resultsIterator.HasNext() {
		// Format query response
		b, status = constructQueryResponseFromIterator(resultsIterator)
		if status != "SUCCESS" {
			return "ERROR", status
		}
	} else {
		return "ERROR", "Error: Query results iterator is empty. Check if function parameters are valid!"
	}

	fmt.Printf("- Ending query GetQueryResultForQueryString queryResult:\n%s\n", b.String())
	return b.String(), "SUCCESS"
}
