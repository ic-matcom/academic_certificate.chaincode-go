package common

import (
	lus "academic_certificates/libutils"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// ContractCommon smart Contract that defines the business logic common
type ContractCommon struct {
	contractapi.Contract
}

// QueryAssetsBy uses a query string to perform a query for any contract asset
// Query string matching state database syntax is passed in and executed as is.
// Supports ad hoc queries that can be defined at runtime by the client.
// Param Ex: {"selector":{"docType":"","id":"myID"}}
//
// Arguments:
//		0: queryStruct map[string]interface{}
// Returns:
//		0: []string
func (cc *ContractCommon) QueryAssetsBy(ctx contractapi.TransactionContextInterface, query map[string]interface{}) ([]interface{}, error) {
	queryString, err := json.MarshalToString(&query)
	if err != nil {
		return nil, err
	}
	return lus.GetQueryResultForQueryString(ctx, queryString)
}

// QueryAssetsWithPagination uses a query string, page size and a bookmark to perform a query
// for assets. Query string matching state database syntax is passed in and executed as is.
// The number of fetched records would be equal to or lesser than the specified page size.
// Supports ad hoc queries that can be defined at runtime by the client.
// If this is not desired, follow the QueryAssetsForOwner example for parameterized queries.
// Only available on state databases that support rich query (e.g. CouchDB)
// Paginated queries are only valid for read only transactions.
// Example: Pagination with Ad hoc Rich Query
func (cc *ContractCommon) QueryAssetsWithPagination(ctx contractapi.TransactionContextInterface, request lus.RichQuerySelector) (*lus.PaginatedQueryResponse, error) {
	queryString, err := json.MarshalToString(request.QueryString)
	if err != nil {
		return nil, err
	} else if queryString == "" {
		return nil, fmt.Errorf("missing query string")
	}

	return lus.GetQueryResultForQueryStringWithPagination(ctx, queryString, int32(request.PageSize), request.Bookmark)
}

func (cc *ContractCommon) GetEvaluateTransactions() []string {
	return []string{"QueryAssetsBy", "QueryAssetsWithPagination", "GetHistory"}
}
