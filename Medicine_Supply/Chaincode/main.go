package main

import (
	"log"

	"medicine/contracts"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	medContract := new(contracts.MedContract)

	chaincode, err := contractapi.NewChaincode(medContract)

	if err != nil {
		log.Panicf("Could not create chaincode : %v", err)
	}

	err = chaincode.Start()

	if err != nil {
		log.Panicf("Failed to start chaincode : %v", err)
	}
}