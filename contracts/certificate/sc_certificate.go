package certificate

import (
	"fmt"

	lus "academic_certificates/libutils"
	"encoding/json"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// ContractCertificate provides functions for managing an asset
type ContractCertificate struct {
	contractapi.Contract
}

type StateValidation uint

const (
	Invalid StateValidation = iota
	One
	Two
	Three
)

//Auxiliary Functions
func (state StateValidation) String() string {
	names := []string{"Invalid", "One", "Two", "Three"}
	if state < Invalid || state > Three {
		return "unknown"
	}
	return names[state]
}

// Asset describes basic details of what makes up a simple asset
type Asset struct {
	DocType       string          `json:"docType"`
	ID            string          `json:"ID"`
	Color         string          `json:"color"`
	Size          int             `json:"size"`
	Owner         string          `json:"owner"`
	OperationType StateValidation `json:"operationType"`
}

type GetRequest struct {
	ID string `json:"id"`
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *Asset
}

// InitLedger adds a base set of cars to the ledger
func (s *ContractCertificate) InitLedger(ctx contractapi.TransactionContextInterface) error {
	assets := []Asset{
		{DocType: lus.CodCert, ID: "asset1", Color: "blue", Size: 5, Owner: "Tomoko", OperationType: One},
		{DocType: lus.CodCert, ID: "asset2", Color: "red", Size: 5, Owner: "Brad", OperationType: Invalid},
		{DocType: lus.CodCert, ID: "asset3", Color: "green", Size: 10, Owner: "Jin Soo", OperationType: Three},
		{DocType: lus.CodCert, ID: "asset4", Color: "yellow", Size: 10, Owner: "Max", OperationType: One},
		{DocType: lus.CodCert, ID: "asset5", Color: "black", Size: 15, Owner: "Adriana", OperationType: Two},
		{DocType: lus.CodCert, ID: "asset6", Color: "white", Size: 15, Owner: "Michel", OperationType: One},
	}

	for i, asset := range assets {
		key, err := ctx.GetStub().CreateCompositeKey(lus.CodCert, []string{"2022", "11", "22", "10302", string(i)})
		if err != nil {
			return err
		}

		assetJSON, err := json.Marshal(asset)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(key, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state: %v", err)
		}
	}

	return nil
}

// CreateAsset issues a new asset to the world state with given details.
func (s *ContractCertificate) CreateAsset(ctx contractapi.TransactionContextInterface, id, color string, size int, owner string) error {
	composeKey, cert, err := lus.ExistsAssetFromId(ctx.GetStub(), lus.CodCert, id)
	if err != nil {
		return err
	} else if cert != nil {
		return fmt.Errorf(lus.ErrorAlreadyExistInState, id)
	}

	asset := Asset{
		DocType:       lus.CodCert,
		ID:            id,
		Color:         color,
		Size:          size,
		Owner:         owner,
		OperationType: One,
	}

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(composeKey, assetJSON)
}

// CreateAssetO issues a new asset to the world state with given details.
func (s *ContractCertificate) CreateAssetO(ctx contractapi.TransactionContextInterface, request *Asset) error {
	composeKey, cert, err := lus.ExistsAssetFromId(ctx.GetStub(), lus.CodCert, request.ID)
	if err != nil {
		return err
	} else if cert != nil {
		return fmt.Errorf(lus.ErrorAlreadyExistInState, request.ID)
	}

	asset := Asset{
		DocType:       lus.CodCert,
		ID:            request.ID,
		Color:         request.Color,
		Size:          request.Size,
		Owner:         request.Owner,
		OperationType: request.OperationType,
	}

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(composeKey, assetJSON)
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *ContractCertificate) ReadAsset(ctx contractapi.TransactionContextInterface, request GetRequest) (*Asset, error) {
	_, assetJSON, err := lus.ExistsAssetFromId(ctx.GetStub(), lus.CodCert, request.ID)
	if err != nil {
		return nil, err
	} else if assetJSON == nil {
		return nil, fmt.Errorf(lus.ErrorNotExistInState, request.ID)
	}

	var asset Asset
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *ContractCertificate) UpdateAsset(ctx contractapi.TransactionContextInterface, request *Asset) error {
	composeKey, assetJSON, err := lus.ExistsAssetFromId(ctx.GetStub(), lus.CodCert, request.ID)
	if err != nil {
		return err
	} else if assetJSON == nil {
		return fmt.Errorf(lus.ErrorNotExistInState, request.ID)
	}

	// overwritting original asset with new asset
	asset := Asset{
		DocType:       lus.CodCert,
		ID:            request.ID,
		Color:         request.Color,
		Size:          request.Size,
		Owner:         request.Owner,
		OperationType: request.OperationType,
	}

	assetJSON, err = json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(composeKey, assetJSON)
}

// DeleteAsset deletes an given asset from the world state.
func (s *ContractCertificate) DeleteAsset(ctx contractapi.TransactionContextInterface, request GetRequest) error {
	composeKey, assetJSON, err := lus.ExistsAssetFromId(ctx.GetStub(), lus.CodCert, request.ID)
	if err != nil {
		return err
	} else if assetJSON == nil {
		return fmt.Errorf(lus.ErrorNotExistInState, request.ID)
	}
	return ctx.GetStub().DelState(composeKey)
}

// TransferAsset updates the owner field of asset with given id in world state, and returns the old owner.
func (s *ContractCertificate) TransferAsset(ctx contractapi.TransactionContextInterface, id string, newOwner string) (string, error) {
	asset, err := s.ReadAsset(ctx, GetRequest{ID: id})
	if err != nil {
		return "", err
	}

	oldOwner := asset.Owner
	asset.Owner = newOwner

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(id, assetJSON)
	if err != nil {
		return "", err
	}

	return oldOwner, nil
}

// GetAllAssets returns all assets found in world state
func (s *ContractCertificate) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	// range query with empty string for startKey and endKey does an open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var results []QueryResult

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		var asset Asset
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}

		queryResult := QueryResult{Key: queryResponse.Key, Record: &asset}
		results = append(results, queryResult)
	}

	return results, nil
}

func (s *ContractCertificate) GetEvaluateTransactions() []string {
	return []string{"CreateAssetO"}
}
