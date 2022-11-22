package lib_utils

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// PaginatedQueryResponse structure used for returning paginated query results and metadata
type PaginatedQueryResponse struct {
	Records             []interface{} `json:"records"` // fabric-contract-api-go not support return []*interfaces{} type
	FetchedRecordsCount int32         `json:"fetchedRecordsCount"`
	Bookmark            string        `json:"bookmark"`
}

// GetQueryResultForQueryString executes the passed in query string.
// The result set is built and returned as a byte array containing the JSON results.
func GetQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]interface{}, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return ConstructQueryResponseFromIterator(resultsIterator)
}

// ConstructQueryResponseFromIterator constructs a slice of assets from the resultsIterator
func ConstructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) ([]interface{}, error) {
	var assets = make([]interface{}, 0)
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var asset interface{}
		err = json.Unmarshal(queryResult.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, asset)
	}
	return assets, nil
}

// GetQueryResultForQueryStringWithPagination executes the passed in query string with
// pagination info. The result set is built and returned as a byte array containing the JSON results.
func GetQueryResultForQueryStringWithPagination(ctx contractapi.TransactionContextInterface, queryString string, pageSize int32, bookmark string) (*PaginatedQueryResponse, error) {
	resultsIterator, responseMetadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, pageSize, bookmark)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	assets, err := ConstructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}

	return &PaginatedQueryResponse{
		Records:             assets,
		FetchedRecordsCount: responseMetadata.FetchedRecordsCount,
		Bookmark:            responseMetadata.Bookmark,
	}, nil
}

// GetState returns state from world state. Unmarshalls the JSON
func GetState(stub shim.ChaincodeStubInterface, key string) ([]byte, error) {
	assetByte, err := stub.GetState(key)
	if err != nil {
		return nil, fmt.Errorf(ErrorWorldState, key)
	} else if assetByte == nil {
		return nil, fmt.Errorf(ErrorNotExistInState, key)
	}

	return assetByte, nil
}

// ExistsAssetFromId returns the compositeKey and value asset if exists in the state database.
//
// Arguments:
//		0: objectType - the target object Type
//		1: assetID - the target ID
// Returns:
//		0: compositeKey
//		1: value of the specified `compositeKey` from the ledger
//		2: error
//
// Possible returns:
//		("", nil, err). Error creating composite key
//		(compositeKey, nil, err). Unable to interact with world state. Asset does not exist in ledger
//		(compositeKey, nil , nil). Asset does not exist in ledger.
//		(compositeKey, assetByte, nil). Asset exists in ledger.
func ExistsAssetFromId(stub shim.ChaincodeStubInterface, objectType string, assetID string) (string, []byte, error) {
	compositeKey, err := CompositeKeyFromID(stub, objectType, assetID)
	if err != nil {
		return "", nil, err
	}
	AssetAsBytes, err := stub.GetState(compositeKey)
	if err != nil {
		return compositeKey, nil, fmt.Errorf("failed to get asset in the world state: %v", err)
	}
	if AssetAsBytes == nil {
		return compositeKey, nil, nil
	}
	return compositeKey, AssetAsBytes, nil
}
