package common

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	lus "github.com/ic-matcom/academic_certificate.chaincode-go/libutils"
)

func (cio *ContractCommon) GetHistory(ctx contractapi.TransactionContextInterface, request *lus.GetHistoryRequest) (lus.HistoryQueryResponse, error) {
	response := lus.HistoryQueryResponse{Response: make([]lus.HistoryAssetPayload, 0)}
	keyAsset, err := lus.CompositeKeyFromID(ctx.GetStub(), request.DocType, request.ID)
	if err != nil {
		return response, err
	}
	resultsIterator, err := ctx.GetStub().GetHistoryForKey(keyAsset)
	if err != nil {
		return response, err
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		responseIterator, err := resultsIterator.Next()
		if err != nil {
			return response, err
		}

		timestamp := lus.GetTimestampRFC3339(responseIterator.Timestamp)
		record := lus.HistoryAssetPayload{
			TxID:  responseIterator.TxId,
			Time:  timestamp,
			Asset: make(map[string]interface{}),
		}

		// if it was not delete operation on given key, then we need to set the
		// corresponding value.
		if !responseIterator.IsDelete {
			var asset map[string]interface{}
			err = json.Unmarshal(responseIterator.Value, &asset)
			if err != nil {
				return response, err
			}
			record.Asset = asset
		}

		response.Response = append(response.Response, record)
	}
	return response, nil
}
