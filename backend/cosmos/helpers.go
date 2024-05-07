package cosmos

import (
	"strings"

	"github.com/AustrianDataLAB/GeWoScout/backend/models"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

func GetQueryItemsPager(container *azcosmos.ContainerClient, city string, query *models.Query) *runtime.Pager[azcosmos.QueryItemsResponse] {
	var sb strings.Builder
	sb.WriteString("SELECT * FROM c WHERE c._partitionKey = @city")

	queryParams := []azcosmos.QueryParameter{
		{Name: "@city", Value: city},
	}

	if query.MinSize != nil && query.MaxSize != nil {
		queryParams = append(queryParams, azcosmos.QueryParameter{Name: "@minSize", Value: *query.MinSize})
		queryParams = append(queryParams, azcosmos.QueryParameter{Name: "@maxSize", Value: *query.MaxSize})
		sb.WriteString(" AND (c.squareMeters BETWEEN @minSize AND @maxSize)")
	} else if query.MinSize != nil {
		queryParams = append(queryParams, azcosmos.QueryParameter{Name: "@minSize", Value: *query.MinSize})
		sb.WriteString(" AND c.squareMeters >= @minSize")
	} else if query.MaxSize != nil {
		queryParams = append(queryParams, azcosmos.QueryParameter{Name: "@maxSize", Value: *query.MaxSize})
		sb.WriteString(" AND c.squareMeters <= @maxSize")
	}

	partitionKey := azcosmos.NewPartitionKeyString(strings.ToLower(city))

	options := azcosmos.QueryOptions{
		QueryParameters:   queryParams,
		ContinuationToken: query.ContinuationToken,
	}

	if query.PageSize != nil {
		options.PageSizeHint = int32(*query.PageSize)
	} else {
		options.PageSizeHint = 10
	}

	return container.NewQueryItemsPager(sb.String(), partitionKey, &options)
}
