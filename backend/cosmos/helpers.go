package cosmos

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/AustrianDataLAB/GeWoScout/backend/models"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

const DEFAULT_PAGE_SIZE = 10

func addQueryParam(sb *strings.Builder, params *[]azcosmos.QueryParameter, name, condition string, value interface{}) {
	*params = append(*params, azcosmos.QueryParameter{Name: name, Value: value})
	sb.WriteString(condition)
}

func GetQueryItemsPager(container *azcosmos.ContainerClient, city string, query *models.Query) *runtime.Pager[azcosmos.QueryItemsResponse] {
	var sb strings.Builder
	sb.WriteString("SELECT * FROM c WHERE c._partitionKey = @city")

	queryParams := []azcosmos.QueryParameter{
		{Name: "@city", Value: city},
	}

	bc := map[models.EnergyClass][]string{
		models.EnergyClassAplusplus: {"A++"},
		models.EnergyClassAplus:     {"A++", "A+"},
		models.EnergyClassA:         {"A++", "A+", "A"},
		models.EnergyClassB:         {"A++", "A+", "A", "B"},
		models.EnergyClassC:         {"A++", "A+", "A", "B", "C"},
		models.EnergyClassD:         {"A++", "A+", "A", "B", "C", "D"},
		models.EnergyClassE:         {"A++", "A+", "A", "B", "C", "D", "E"},
	}

	fieldMappings := map[string]struct {
		condition string
		value     interface{}
	}{
		"title":                {" AND c.title = @title", query.Title},
		"housingCooperative":   {" AND c.housingCooperative = @housingCooperative", query.HousingCooperative},
		"projectId":            {" AND c.projectId = @projectId", query.ProjectId},
		"postalCodes":          {" AND ARRAY_CONTAINS(@postalCodes, c.postalCode) = true", query.PostalCode},
		"roomCount":            {" AND c.roomCount = @roomCount", query.RoomCount},
		"minRoomCount":         {" AND c.roomCount >= @minRoomCount", query.MinRoomCount},
		"maxRoomCount":         {" AND c.roomCount <= @maxRoomCount", query.MaxRoomCount},
		"minSquareMeters":      {" AND c.squareMeters >= @minSquareMeters", query.MinSquareMeters},
		"maxSquareMeters":      {" AND c.squareMeters <= @maxSquareMeters", query.MaxSquareMeters},
		"availableFrom":        {" AND c.availabilityDate <= @availableFrom", query.AvailableFrom},
		"minYearBuilt":         {" AND c.yearBuilt >= @minYearBuilt", query.MinYearBuilt},
		"maxYearBuilt":         {" AND c.yearBuilt >= @maxYearBuilt", query.MaxYearBuilt},
		"minHwgEnergyClass":    {" AND ARRAY_CONTAINS(@minHwgEnergyClass, c.hwgEnergyClass) = true", query.MinHwgEnergyClass},
		"minFgeeEnergyClass":   {" AND ARRAY_CONTAINS(@minFgeeEnergyClass, c.fgeeEnergyClass) = true", query.MinFgeeEnergyClass},
		"listingType":          {" AND c.listingType = @listingType", query.ListingType},
		"minRentPricePerMonth": {" AND c.rentPrice >= @minRent", query.MinRentPricePerMonth},
		"maxRentPricePerMonth": {" AND c.rentPrice <= @maxRent", query.MaxRentPricePerMonth},
		"minCooperativeShare":  {" AND c.cooperativeShare >= @minCooperativeShare", query.MinCooperativeShare},
		"maxCooperativeShare":  {" AND c.cooperativeShare <= @maxCooperativeShare", query.MaxCooperativeShare},
		"minSalePrice":         {" AND c.salePrice >= @minSalePrice", query.MinSalePrice},
		"maxSalePrice":         {" AND c.salePrice >= @maxSalePrice", query.MaxSalePrice},
	}

	for field, mapping := range fieldMappings {
		if !reflect.ValueOf(mapping.value).IsNil() {
			if field == "minHwgEnergyClass" || field == "minFgeeEnergyClass" {
				ecStr, _ := (mapping.value).(*models.EnergyClass)
				if *ecStr != "F" {
					addQueryParam(&sb, &queryParams, "@"+field, mapping.condition, bc[*ecStr])
				}
			} else {
				addQueryParam(&sb, &queryParams, "@"+field, mapping.condition, mapping.value)
			}
		}
	}

	if query.SortBy != nil {
		sb.WriteString(fmt.Sprintf(" ORDER BY c.%s ", *query.SortBy))
		if query.SortType != nil {
			sb.WriteString(string(*query.SortType))
		} else {
			sb.WriteString("ASC")
		}
	}

	partitionKey := azcosmos.NewPartitionKeyString(strings.ToLower(city))

	options := azcosmos.QueryOptions{
		QueryParameters:   queryParams,
		ContinuationToken: query.ContinuationToken,
	}

	if query.PageSize != nil {
		options.PageSizeHint = int32(*query.PageSize)
	} else {
		options.PageSizeHint = DEFAULT_PAGE_SIZE
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
