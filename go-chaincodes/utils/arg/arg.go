/*
Package arg provides an auxiliary function that receives arguments
and manipulates them to be usefull on cross-chaincode comunication.
*/
package arg

// ToChaincodeArgs - prepares function arguments to invoke
// params: args
func ToChaincodeArgs(args ...string) [][]byte {
	bargs := make([][]byte, len(args))
	for i, arg := range args {
		bargs[i] = []byte(arg)
	}
	return bargs
}
