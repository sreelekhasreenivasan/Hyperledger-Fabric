package contracts

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type MedContract struct{
	contractapi.Contract
}

type Medicine struct {

	BatchId           string `json:"batchId"`
	Name              string `json:"name"`
	Quantity		  int    `json:"quantity"`
	Manufacturer      string `json:"manufacturer"`
	DateOfManufacture string `json:"dateOfManufacture"`
	ExpiringDate      string `json:"expiringDate"`
	CurrentHolder     string `json:"currentHolder"`
	Status            string `json:"status"`

}

//product already exist
func (m * MedContract) ExistProduct(ctx contractapi.TransactionContextInterface, batchId string)(bool,error){
	result,err :=ctx.GetStub().GetState(batchId)

	if err!=nil{
		return false, fmt.Errorf("failed to read from world state:%v",err)
	}
	return result!=nil,nil
}

//create product
func (m * MedContract) CreateProduct(ctx contractapi.TransactionContextInterface, 
	batchId string,
	name string,
	quantity int,
	manufacturer string,
	dateOfManufacture string,
	expiringDate string)(string,error){

	manufacturerId,err := ctx.GetClientIdentity().GetMSPID()

	if err!=nil{
		return "",err
	}

	if manufacturerId == "manufacturer-med-com"{

		existProduct,err := m.ExistProduct(ctx,batchId)
		
		if err!=nil{
			return "", fmt.Errorf("failed to load data %s", err)
		} else if existProduct{
			return "", fmt.Errorf("the product, %s already exists", batchId)
		}

		product := Medicine{
			BatchId: batchId,
			Name: name,
			Quantity: quantity,
			Manufacturer: manufacturer,
			DateOfManufacture: dateOfManufacture,
			ExpiringDate: expiringDate,
			CurrentHolder: "In Factory",
			Status: "On Progress",
		}
		bytes,_ :=json.Marshal(product)

		err = ctx.GetStub().PutState(batchId, bytes)

		if err!= nil{

			return "", err

		} else{

			return fmt.Sprintf("successfully added product %v", batchId), nil
		}

	} else{

		return "", fmt.Errorf("user under following MSPID: %v can't perform this action", batchId)
	}

}


//get product details
func(m* MedContract) ReadProdcut(ctx contractapi.TransactionContextInterface,batchId string)(*Medicine,error){

	bytes,err:=ctx.GetStub().GetState(batchId)

	if err!=nil{
		return nil, fmt.Errorf("failed to read data from world state:%v",err)
	}
	if bytes==nil{
		return nil, fmt.Errorf("the product %s does not exist",batchId)
	}

	var product Medicine

	err = json.Unmarshal(bytes,&product)

	if err != nil {
		return nil, fmt.Errorf("could not unmarshal world state data to type Car")
	}
	return &product,nil

}

//update product details
func(m*MedContract) UpdateProduct(ctx contractapi.TransactionContextInterface,
	batchId string,
	name string,
	quantity int,
	currentHolder string,
	status string,)(string,error){

		existProduct,err := m.ExistProduct(ctx,batchId)

		if err!=nil{
			return "", fmt.Errorf("%s",err)
		}
		if !existProduct{
			return "", fmt.Errorf("product %s does not exist",batchId)
		}

		product,err:=m.ReadProdcut(ctx,batchId)

		if err!=nil{
			return "",err
		}

		product.Name = name
		product.Quantity = quantity
		product.CurrentHolder = currentHolder
		product.Status = status

		bytes, err := json.Marshal(product)
		if err != nil {
		return "", err
		}

		err = ctx.GetStub().PutState(batchId, bytes)
		if err != nil {
		return "", err
		}

		return fmt.Sprintf("product %s successfully updated", batchId), nil


	}

	//delete product
	func (m *MedContract) DeleteProduct(ctx contractapi.TransactionContextInterface, batchId string) (string, error) {

		manufacturerId, err := ctx.GetClientIdentity().GetMSPID()
		if err != nil {
			return "", err
		}
		if manufacturerId == "manufacturer-med-com" {
	
			existProduct, err := m.ExistProduct(ctx, batchId)
			if err != nil {
				return "", fmt.Errorf("%s", err)
			} else if !existProduct {
				return "", fmt.Errorf("the product %s does not exist", batchId)
			}
	
			err = ctx.GetStub().DelState(batchId)
			if err != nil {
				return "", err
			} else {
				return fmt.Sprintf("product with batchId %v is deleted from the world state.", batchId), nil
			}
	
		} else {
			return "", fmt.Errorf("user under following MSPID: %v can't perform this action", manufacturerId)
		}
	}

	//transfer product
	func(m*MedContract) TransferProduct(ctx contractapi.TransactionContextInterface,batchId string, newHolder string)(string,error){

		existProduct,err := m.ExistProduct(ctx,batchId)

		if err!=nil{
			return "", fmt.Errorf("%s",err)
		}
		if !existProduct{
			return "", fmt.Errorf("product %s does not exist",batchId)
		}

		product,err:=m.ReadProdcut(ctx,batchId)

		if err!=nil{
			return "",err
		}

		product.CurrentHolder = newHolder
		product.Status = "In transit"

		bytes, err := json.Marshal(product)
		if err != nil {
		return "", err
		}

		err = ctx.GetStub().PutState(batchId, bytes)
		if err != nil {
		return "", err
		}

		return fmt.Sprintf("product %s transferred to %s", batchId, newHolder), nil

	}
