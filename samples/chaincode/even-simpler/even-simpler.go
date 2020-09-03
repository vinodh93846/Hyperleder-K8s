package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	sc "github.com/hyperledger/fabric-protos-go/peer"
	//  "github.com/hyperledger/fabric/core/chaincode/shim"
	//  sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the Bar structure, with 4 properties.  Structure tags are used by encoding/json library
type Bar struct {
	BarID              string `json:"id"`
	BarLocation        string `json:"barLocation"`
	BarSerialNumber    string `json:"barSerialNumber"`
	Purity             string `json:"purity"`
	BarRefiner         string `json:"barRefiner"`
	BarHallmarkVerfied string `json:"barHallmarkVerfied"`
	BarWeightInGms     string `json:"barWeightInGms"`
}
type Buy struct {
	DSGID          string `json:"id"`
	OrderID        string `json:"orderId"`
	Amount         string `json:"amount"`
	AmountWithFees string `json:"amountWithFees"`
	Stage          string `json:"stage"`
	PaymentStatus  string `json:"paymentStatus"`
	EstimatedGrams string `json:"estimatedGrams"`
	UserID         string `json:"userId"`
}
type Sell struct {
	DSGID           string `json:"id"`
	OrderID         string `json:"orderId"`
	Grams           string `json:"grams"`
	EstimatedAmount string `json:"estimatedamount"`
	UserID          string `json:"userId"`
}
type Send struct {
	DSGID          string `json:"id"`
	OrderID        string `json:"orderId"`
	Grams          string `json:"grams"`
	SenderUserID   string `json:"senderUserId"`
	ReceiverUserID string `json:"receiverUserId"`
}
type Trade struct {
	DSGID   string `json:"id"`
	OrderID string `json:"orderId"`
	Grams   string `json:"grams"`
	UserID  string `json:"userId"`
}

func GetUID() (string, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	return id.String(), err
}

/*
 * The Init method is called when the Smart Contract "DSG Bar" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "DSG Bar"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "getBar" {
		return s.getBar(APIstub, args)
	} else if function == "createBar" {
		return s.createBar(APIstub, args)
	} else if function == "queryAllBars" {
		return s.queryAllBars(APIstub)
	} else if function == "getBuy" {
		return s.getBuy(APIstub, args)
	} else if function == "createBuy" {
		return s.createBuy(APIstub, args)
	} else if function == "queryAllBuys" {
		return s.queryAllBuys(APIstub)
	} else if function == "getSell" {
		return s.getSell(APIstub, args)
	} else if function == "createSell" {
		return s.createSell(APIstub, args)
	} else if function == "queryAllSells" {
		return s.queryAllSells(APIstub)
	} else if function == "getSend" {
		return s.getSend(APIstub, args)
	} else if function == "createSend" {
		return s.createSend(APIstub, args)
	} else if function == "queryAllSends" {
		return s.queryAllSends(APIstub)
	} else if function == "getTrade" {
		return s.getTrade(APIstub, args)
	} else if function == "createTrade" {
		return s.createTrade(APIstub, args)
	} else if function == "queryAllTreades" {
		return s.queryAllTrades(APIstub)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) getBar(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	barAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(barAsBytes)
}
func (s *SmartContract) createBar(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	fmt.Printf("Adding Bar to the ledger ...\n")
	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
	}
	uid, err := GetUID()
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	id := "Bar-" + uid
	fmt.Printf("Validating Bar data\n")
	//Validate the Org data
	var bar = Bar{BarID: id,
		BarLocation:        args[0],
		BarSerialNumber:    args[1],
		Purity:             args[2],
		BarRefiner:         args[3],
		BarHallmarkVerfied: args[4],
		BarWeightInGms:     args[5]}

	barAsBytes, _ := json.Marshal(bar)
	APIstub.PutState(args[0], barAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) queryAllBars(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "Bar-000"
	endKey := "Bar-9999999999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllBars:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}
func (s *SmartContract) getBuy(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	buyAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(buyAsBytes)
}
func (s *SmartContract) createBuy(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	fmt.Printf("Adding Buy to the ledger ...\n")
	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 8")
	}
	uid, err := GetUID()
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	id := "DSG-" + uid
	fmt.Printf("Validating Bar data\n")
	//Validate the Org data
	var buy = Buy{DSGID: id,
		OrderID:        args[0],
		Amount:         args[1],
		AmountWithFees: args[2],
		Stage:          args[3],
		PaymentStatus:  args[4],
		EstimatedGrams: args[5],
		UserID:         args[6]}

	buyAsBytes, _ := json.Marshal(buy)
	APIstub.PutState(args[0], buyAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) queryAllBuys(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "Buy-"
	endKey := "Buy-9999999999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllBuys:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}
func (s *SmartContract) getSell(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	sellAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(sellAsBytes)
}
func (s *SmartContract) createSell(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	fmt.Printf("Adding Sell to the ledger ...\n")
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
	}
	uid, err := GetUID()
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	id := "DSG-" + uid
	fmt.Printf("Validating Bar data\n")
	//Validate the Org data
	var sell = Sell{DSGID: id,
		OrderID:         args[0],
		Grams:           args[1],
		EstimatedAmount: args[2],
		UserID:          args[3]}

	sellAsBytes, _ := json.Marshal(sell)
	APIstub.PutState(args[0], sellAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) queryAllSells(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "DSG-000"
	endKey := "DSG-9999999999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllSells:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}
func (s *SmartContract) getSend(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	sendAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(sendAsBytes)
}
func (s *SmartContract) createSend(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	fmt.Printf("Adding Send to the ledger ...\n")
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
	}
	uid, err := GetUID()
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	id := "DSG-" + uid
	fmt.Printf("Validating Bar data\n")
	//Validate the Org data
	var send = Send{DSGID: id,
		OrderID:        args[0],
		Grams:          args[1],
		SenderUserID:   args[2],
		ReceiverUserID: args[3]}

	sendAsBytes, _ := json.Marshal(send)
	APIstub.PutState(args[0], sendAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) queryAllSends(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "DSG-000"
	endKey := "DSG-9999999999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllSends:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}
func (s *SmartContract) getTrade(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	tradeAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(tradeAsBytes)
}
func (s *SmartContract) createTrade(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	fmt.Printf("Adding Trade to the ledger ...\n")
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
	}
	uid, err := GetUID()
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	id := "DSG-" + uid
	fmt.Printf("Validating Bar data\n")
	//Validate the Org data
	var trade = Trade{DSGID: id,
		OrderID: args[0],
		Grams:   args[1],
		UserID:  args[2]}

	tradeAsBytes, _ := json.Marshal(trade)
	APIstub.PutState(args[0], tradeAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) queryAllTrades(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "DSG-000"
	endKey := "DSG-9999999999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllTrades:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
