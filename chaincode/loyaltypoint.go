package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {

	/*var err error

	valAsbytes := [...]byte{'1', '0', '0', '0'}
	pointString := string(valAsbytes[:])
	points, err := strconv.ParseInt(pointString, 10, 64)
	if err != nil {

	}

	points = points - 30

	fmt.Printf("pointString : %s, points: %d", pointString, points)*/

	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s .", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var mCardID, lPointString, exchangeStoreNameString, exRateString string
	var err error
	fmt.Println("running initialData()")

	if len(args) != 4 {
		return nil, errors.New("incorrect number of arguments. Expecting 4. M-card id, point, store and exchange rate")
	}

	mCardID = args[0] //rename for fun
	lPointString = args[1]
	exchangeStoreNameString = args[2]
	exRateString = args[3]

	//write the variable into the chaincode state
	err = stub.PutState(mCardID, []byte(lPointString))
	if err != nil {
		return nil, err
	}

	err = stub.PutState("exRate_"+exchangeStoreNameString, []byte(exRateString))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" { //initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else if function == "redeemPoint" {
		return t.redeemPoint(stub, args)
	} else if function == "setNewRate" {
		return t.setNewRate(stub, args)
	} else if function == "addMCard" {
		return t.addMCard(stub, args)
	} else if function == "deleteAllState" {
		return t.deleteAllState(stub)
	}

	fmt.Println("invoke did not find func: " + function) //error

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read_point" {
		return t.readPoint(stub, args)
	} else if function == "read_exchangeRate" {
		return t.readExchangeRate(stub, args)
	}

	fmt.Println("query did not find func: " + function) //error

	return nil, errors.New("Received unknown function query: " + function)
}

func (t *SimpleChaincode) readPoint(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting M-card id to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}

func (t *SimpleChaincode) readExchangeRate(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of store to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState("exRate_" + key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}

func (t *SimpleChaincode) redeemPoint(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running redeemPoint()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. M-card id and redeem points to set")
	}

	key = args[0] //rename for fun
	value = args[1]

	valAsbytes, err := stub.GetState(key)
	if err != nil {
		var jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	// Convert to int
	pointString := string(valAsbytes[:])
	points, err := strconv.ParseInt(pointString, 10, 64)
	if err != nil {
		return nil, err
	}

	redeemPoint, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return nil, err
	}

	// Redeem points
	points -= redeemPoint

	// To string
	remainPointsString := strconv.FormatInt(points, 10)

	err = stub.PutState(key, []byte(remainPointsString)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (t *SimpleChaincode) setNewRate(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running setNewRate()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of store and exchange rate to set")
	}

	key = args[0] //rename for fun
	value = args[1]
	err = stub.PutState("exRate_"+key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (t *SimpleChaincode) addMCard(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running addMCard()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. m-card id and points to set")
	}

	key = args[0] //rename for fun
	value = args[1]
	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (t *SimpleChaincode) deleteAllState(stub shim.ChaincodeStubInterface) ([]byte, error) {
	var err error
	fmt.Println("running deleteAllState()")

	iter, err := stub.RangeQueryState("", "")

	if err != nil {
		fmt.Printf("Error deleting table: %s", err)
		return nil, err
	}
	defer iter.Close()
	for iter.HasNext() {
		key, _, err := iter.Next()
		if err != nil {
			fmt.Printf("Error deleting table: %s", err)
			return nil, err
		}
		err = stub.DelState(key)
		if err != nil {
			fmt.Printf("Error deleting table: %s", err)
			return nil, err
		}
	}

	return nil, nil
}
