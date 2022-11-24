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
	var assets []Asset
	for i := 0; i < 10; i++ {
		assets = append(assets, Asset{
			DocType:               lus.CodCert,
			Certification:         "Licenciado en Derecho",
			GoldCertificate:       false,
			Emitter:               "Universidad de La Habana",
			Accredited:            fmt.Sprintf("Joe Doe %d", i),
			Date:                  "8 de Noviembre del 2010",
			SecretaryValidating:   "Mirtha Guerra",
			DeanValidating:        "",
			RectorValidating:      "",
			FacultyVolumeFolio:    "254,136",
			UniversityVolumeFolio: "158,187",
			InvalidReason:         "",
			Status:                SignedS,
		})
	}
	for i := 10; i < 20; i++ {
		assets = append(assets, Asset{
			DocType:               lus.CodCert,
			Certification:         "Licenciado en Química",
			GoldCertificate:       true,
			Emitter:               "Universidad de La Habana",
			Accredited:            fmt.Sprintf("Joe Doe %d", i),
			Date:                  "10 de Julio del 2018",
			SecretaryValidating:   "Manuela Azurra",
			DeanValidating:        "Pedro Navaja",
			RectorValidating:      "",
			FacultyVolumeFolio:    "254, 333",
			UniversityVolumeFolio: "158,781",
			InvalidReason:         "",
			Status:                SignedSD,
		})
	}

	for i, asset := range assets {
		key, err := ctx.GetStub().CreateCompositeKey(lus.CodCert, []string{"2022", "11", "22", "10302", string(rune(i + 1))})
		if err != nil {
			return err
		}
		asset.ID = key

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

// CreateAssetOld issues a new asset to the world state with given details.
func (s *ContractCertificate) CreateAssetOld(ctx contractapi.TransactionContextInterface, id string,
	certification string, goldCertificate bool, emitter string, accredited string, date string, createdBy,
	facultyVolumeFolio string, universityVolumeFolio string) error {
	composeKey, cert, err := lus.ExistsAssetFromId(ctx.GetStub(), lus.CodCert, id)
	if err != nil {
		return err
	} else if cert != nil {
		return fmt.Errorf(lus.ErrorAlreadyExistInState, id)
	}

	asset := Asset{
		DocType:               lus.CodCert,
		ID:                    id,
		Certification:         certification,
		GoldCertificate:       goldCertificate,
		Emitter:               emitter,
		Accredited:            accredited,
		Date:                  date,
		CreatedBy:             createdBy,
		SecretaryValidating:   "",
		DeanValidating:        "",
		RectorValidating:      "",
		FacultyVolumeFolio:    facultyVolumeFolio,
		UniversityVolumeFolio: universityVolumeFolio,
		InvalidReason:         "",
		Status:                New,
	}

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(composeKey, assetJSON)
}

// CreateAsset issues a new asset to the world state with given details.
func (s *ContractCertificate) CreateAsset(ctx contractapi.TransactionContextInterface, request *Asset) error {
	composeKey, cert, err := lus.ExistsAssetFromId(ctx.GetStub(), lus.CodCert, request.ID)
	if err != nil {
		return err
	} else if cert != nil {
		return fmt.Errorf(lus.ErrorAlreadyExistInState, request.ID)
	}

	asset := Asset{
		DocType:               lus.CodCert,
		ID:                    request.ID,
		Certification:         request.Certification,
		GoldCertificate:       request.GoldCertificate,
		Emitter:               request.Emitter,
		Accredited:            request.Accredited,
		Date:                  request.Date,
		CreatedBy:             request.CreatedBy,
		SecretaryValidating:   "",
		DeanValidating:        "",
		RectorValidating:      "",
		FacultyVolumeFolio:    request.FacultyVolumeFolio,
		UniversityVolumeFolio: request.UniversityVolumeFolio,
		InvalidReason:         "",
		Status:                New,
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
	// Check new params of the asset consistency

	// If certificate is valid then it should have the 3 signatures
	if (request.Status == Valid) && (request.SecretaryValidating == "" || request.DeanValidating == "" || request.RectorValidating == "") {
		return fmt.Errorf(lus.ErrorInconsistentStatus)
	}
	// If certificate is SignedSD then it should have Secretary and Dean signatures
	if (request.Status == SignedSD) && (request.SecretaryValidating == "" || request.DeanValidating == "") {
		return fmt.Errorf(lus.ErrorInconsistentStatus)
	}
	// If certificate is SignedSD then it should have Secretary and Dean signatures
	if (request.Status == SignedS) && (request.SecretaryValidating == "") {
		return fmt.Errorf(lus.ErrorInconsistentStatus)
	}
	// If certificate is revoked then it should have a revoked reason
	if (request.Status == Invalid) && (request.InvalidReason == "") {
		return fmt.Errorf(lus.ErrorInconsistentInvalidation)
	}
	// overwritting original asset with new asset
	asset := Asset{
		DocType:               lus.CodCert,
		ID:                    request.ID,
		Certification:         request.Certification,
		GoldCertificate:       request.GoldCertificate,
		Emitter:               request.Emitter,
		Accredited:            request.Accredited,
		Date:                  request.Date,
		CreatedBy:             request.CreatedBy,
		SecretaryValidating:   request.SecretaryValidating,
		DeanValidating:        request.DeanValidating,
		RectorValidating:      request.RectorValidating,
		FacultyVolumeFolio:    request.FacultyVolumeFolio,
		UniversityVolumeFolio: request.UniversityVolumeFolio,
		InvalidReason:         request.InvalidReason,
		Status:                request.Status,
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

func (s *ContractCertificate) GetEvaluateTransactions() []string {
	return []string{"CreateAssetO"}
}
