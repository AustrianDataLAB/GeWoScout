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

func getValidSortParams() map[string]bool {
	return map[string]bool{
		"title":              true,
		"housingCooperative": true,
		"projectId":          true,
		"listingId":          true,
		"country":            true,
		"city":               true,
		"postalCode":         true,
		"address":            true,
		"roomCount":          true,
		"squareMeters":       true,
		"availabilityDate":   true,
		"yearBuilt":          true,
		"hwgEnergyClass":     true,
		"fgeeEnergyClass":    true,
		"listingType":        true,
		"rentPricePerMonth":  true,
		"cooperativeShare":   true,
		"salePrice":          true,
		"additionalFees":     true,
		"detailsUrl":         true,
		"previewImageUrl":    true,
		"scraperId":          true,
		"createdAt":          true,
	}
}

type fieldMapping struct {
	condition string
	value     interface{}
}

func getFieldMappings(query *models.ListingsQuery) map[string]fieldMapping {
	return map[string]fieldMapping{
		"title":                {" AND CONTAINS(LOWER(c.title), LOWER(@title))", query.Title},
		"housingCooperative":   {" AND CONTAINS(LOWER(c.housingCooperative), LOWER(@housingCooperative))", query.HousingCooperative},
		"projectId":            {" AND LOWER(c.projectId) = LOWER(@projectId)", query.ProjectId},
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
}

func GetListingsQueryItemsPager(
	container *azcosmos.ContainerClient,
	city string,
	query *models.ListingsQuery,
) *runtime.Pager[azcosmos.QueryItemsResponse] {
	var sb strings.Builder
	sb.WriteString("SELECT * FROM c WHERE c._partitionKey = @city")

	queryParams := []azcosmos.QueryParameter{
		{Name: "@city", Value: city},
	}

	fieldMappings := getFieldMappings(query)

	for field, mapping := range fieldMappings {
		if !reflect.ValueOf(mapping.value).IsNil() {
			if field == "minHwgEnergyClass" || field == "minFgeeEnergyClass" {
				ecStr, _ := (mapping.value).(*models.EnergyClass)
				if *ecStr != "F" {
					addQueryParam(&sb, &queryParams, "@"+field, mapping.condition, models.GetEnergyClasses()[:ecStr.GetIndex()+1])
				}
			} else if field == "postalCodes" {
				postalCodeStr := mapping.value.(*string)
				addQueryParam(&sb, &queryParams, "@"+field, mapping.condition, strings.Split(*postalCodeStr, ","))
			} else {
				addQueryParam(&sb, &queryParams, "@"+field, mapping.condition, mapping.value)
			}
		}
	}

	if query.SortBy != nil {
		if _, ok := getValidSortParams()[*query.SortBy]; ok {
			sb.WriteString(fmt.Sprintf(" ORDER BY c.%s ", *query.SortBy))
			if query.SortType != nil && *query.SortType == "DESC" {
				sb.WriteString(string(*query.SortType))
			} else {
				sb.WriteString("ASC")
			}
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

func GetUserData(ctx context.Context, container *azcosmos.ContainerClient, userId string) ([]models.UserData, error) {
	queryParams := []azcosmos.QueryParameter{
		{Name: "@userId", Value: userId},
	}
	options := azcosmos.QueryOptions{
		QueryParameters: queryParams,
	}
	pager := container.NewQueryItemsPager("SELECT * FROM c WHERE c._partitionKey = @userId", azcosmos.NewPartitionKeyString(userId), &options)

	uds := []models.UserData{}

	for pager.More() {
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		page, err := pager.NextPage(ctx)
		cancel()
		if err != nil {
			return nil, fmt.Errorf("failed to get next page: %w", err)
		}

		for _, item := range page.Items {
			ud := models.UserData{}
			if err := json.Unmarshal(item, &ud); err != nil {
				return nil, fmt.Errorf("failed to unmarshal item: %w", err)
			}

			ud.City = &ud.Id

			uds = append(uds, ud)
		}
	}

	return uds, nil
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

func getPreferenceFieldMappings(listing *models.Listing) map[string]fieldMapping {
	return map[string]fieldMapping{
		"title":              {" AND c.title ? (CONTAINS(LOWER(@title), LOWER(c.title)) = true) : true", listing.Title},
		"housingCooperative": {" AND c.housingCooperative ? (CONTAINS(LOWER(@housingCooperative), LOWER(c.housingCooperative)) = true) : true", listing.HousingCooperative},
		"projectId":          {" AND c.projectId ? (LOWER(c.projectId) = LOWER(@projectId)) : true", listing.ProjectID},
		"postalCodes":        {" AND c.postalCode ? (CONTAINS(c.postalCode, @postalCodes) = true) : true", listing.PostalCode},
		"roomCount":          {" AND (c.minRoomCount ? (c.minRoomCount <= @roomCount) : true) AND (c.maxRoomCount ? (c.maxRoomCount >= @roomCount) : true)", listing.RoomCount},
		"squareMeters":       {" AND (c.minSqm ? (c.minSqm <= @squareMeters) : true) AND (c.maxSqm ? (c.maxSqm >= @squareMeters) : true)", listing.SquareMeters},
		"availabilityDate":   {" AND c.availableFrom ? (c.availableFrom <= @availabilityDate) : true", listing.AvailabilityDate},
		"yearBuilt":          {" AND (c.minYearBuilt ? (c.minYearBuilt <= @yearBuilt) : true) AND (c.maxYearBuilt ? (c.maxYearBuilt >= @yearBuilt) : true)", listing.YearBuilt},
		"hwgEnergyClass":     {" AND c.minHwgEnergyClass ? (ARRAY_CONTAINS(@hwgEnergyClass, c.minHwgEnergyClass) = true) : true", listing.HwgEnergyClass},
		"fgeeEnergyClass":    {" AND c.minFgeeEnergyClass ? (ARRAY_CONTAINS(@fgeeEnergyClass, c.minFgeeEnergyClass) = true) : true", listing.FgeeEnergyClass},
		"listingType":        {" AND c.listingType ? (c.listingType = @listingType) : true", listing.ListingType},
		"rentPricePerMonth":  {" AND (c.minRentPrice ? (c.minRentPrice <= @rentPricePerMonth) : true) AND (c.maxRentPrice ? (c.maxRentPrice >= @rentPricePerMonth) : true)", listing.RentPricePerMonth},
		"cooperativeShare":   {" AND (c.minCooperativeShare ? (c.minCooperativeShare <= @cooperativeShare) : true) AND (c.maxCooperativeShare ? (c.maxCooperativeShare >= @cooperativeShare) : true)", listing.CooperativeShare},
		"salePrice":          {" AND (c.minSalePrice ? (c.minSalePrice <= @salePrice) : true) AND (c.maxSalePrice ? (c.maxSalePrice >= @salePrice) : true)", listing.SalePrice},
	}
}

func GetUsersMatchingWithListing(ctx context.Context, container *azcosmos.ContainerClient, listing *models.Listing) ([]string, error) {
	var sb strings.Builder
	sb.WriteString("SELECT c.email FROM c WHERE c._partitionKey = @city")

	queryParams := []azcosmos.QueryParameter{
		{Name: "@city", Value: listing.PartitionKey},
	}

	fieldMappings := getPreferenceFieldMappings(listing)

	for field, mapping := range fieldMappings {
		switch mapping.value.(type) {
		case string:
			if mapping.value != "" {
				if field == "hwgEnergyClass" || field == "fgeeEnergyClass" {
					ecStr, ok := (mapping.value).(string)
					if !ok {
						return nil, fmt.Errorf("energy class has incorrect format")
					}
					ecClass := models.EnergyClass(ecStr)
					addQueryParam(&sb, &queryParams, "@"+field, mapping.condition, models.GetEnergyClasses()[ecClass.GetIndex():])
				} else {
					addQueryParam(&sb, &queryParams, "@"+field, mapping.condition, mapping.value)
				}
			}
		case int:
			addQueryParam(&sb, &queryParams, "@"+field, mapping.condition, mapping.value)
		case *int:
			if mapping.value != nil {
				addQueryParam(&sb, &queryParams, "@"+field, mapping.condition, mapping.value)
			}
		default:
			continue
		}

	}

	partitionKey := azcosmos.NewPartitionKeyString(listing.PartitionKey)

	options := azcosmos.QueryOptions{
		QueryParameters: queryParams,
	}

	pager := container.NewQueryItemsPager(sb.String(), partitionKey, &options)

	emails := []string{}

	type EmailItem struct {
		Email string `json:"email"`
	}

	if pager.More() {
		response, err := pager.NextPage(ctx)
		if err != nil {
			return []string{}, fmt.Errorf("failed to get matching user preferences: %s", err.Error())
		}

		for _, bytes := range response.Items {
			var emailItem EmailItem
			if err := json.Unmarshal(bytes, &emailItem); err != nil {
				return []string{}, fmt.Errorf("failed to unmarshal email item: %s", err.Error())
			}
			emails = append(emails, emailItem.Email)
		}
	}

	return emails, nil
}
