package cosmos

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"

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

func GetNonExistingIds(ctx context.Context, container *azcosmos.ContainerClient, idsByPk map[string][]string) (map[string][]string, error) {
	nonExistingIds := make(map[string][]string)

	for pk, ids := range idsByPk {
		nonExistingIds[pk] = []string{}

		existingIds, err := GetExistingIds(ctx, container, map[string][]string{pk: ids})
		if err != nil {
			return nil, err
		}

		for _, id := range ids {
			if !slices.Contains(existingIds[pk], id) {
				nonExistingIds[pk] = append(nonExistingIds[pk], id)
			}
		}
	}

	return nonExistingIds, nil
}

func GetExistingIds(ctx context.Context, container *azcosmos.ContainerClient, idsByPk map[string][]string) (map[string][]string, error) {
	existingIds := make(map[string][]string)

	for pk, ids := range idsByPk {

		existingIds[pk] = []string{}

		partitionKey := azcosmos.NewPartitionKeyString(pk)
		placeholder, parameters := genStringParamList(ids)

		query := fmt.Sprintf("SELECT c.id FROM c WHERE c.id IN (%s)", strings.Join(placeholder, ", "))

		pager := container.NewQueryItemsPager(query, partitionKey, &azcosmos.QueryOptions{
			QueryParameters: parameters,
		})

		for pager.More() {
			ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
			page, err := pager.NextPage(ctx)
			cancel()
			if err != nil {
				return nil, fmt.Errorf("failed to get next page: %w", err)
			}

			type IdItem struct {
				ID string `json:"id"`
			}

			for _, item := range page.Items {
				var idItem IdItem
				err := json.Unmarshal(item, &idItem)
				if err != nil {
					return nil, fmt.Errorf("failed to unmarshal item: %w", err)
				}

				existingIds[pk] = append(existingIds[pk], idItem.ID)
			}
		}
	}

	return existingIds, nil
}

func genStringParamList(ids []string) ([]string, []azcosmos.QueryParameter) {
	placeholders := make([]string, len(ids))
	parameters := make([]azcosmos.QueryParameter, len(ids))
	for i, id := range ids {
		placeholder := fmt.Sprintf("@id%d", i)
		placeholders[i] = placeholder
		parameters[i] = azcosmos.QueryParameter{
			Name:  placeholder,
			Value: id,
		}
	}
	return placeholders, parameters
}
