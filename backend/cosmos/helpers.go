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
			switch mapping.value.(type) {
			case *float32:
				v := mapping.value.(*float32)
				addQueryParam(&sb, &queryParams, "@"+field, mapping.condition, *v)
			default:
				if field == "minHwgEnergyClass" || field == "minFgeeEnergyClass" {
					ecStr, _ := (mapping.value).(*models.EnergyClass)
					if *ecStr != models.EnergyClassG {
						addQueryParam(&sb, &queryParams, "@"+field, mapping.condition, models.GetEnergyClasses()[:ecStr.GetIndex()+1])
					}
				} else if field == "postalCodes" {
					postalCodeStr := mapping.value.(*string)
					addQueryParam(&sb, &queryParams, "@"+field, mapping.condition, strings.Split(*postalCodeStr, ","))
				} else if field == "listingType" {
					listingType := mapping.value.(*models.ListingType)
					if *listingType != "both" {
						addQueryParam(&sb, &queryParams, "@"+field, mapping.condition, *listingType)
					}
				} else {
					addQueryParam(&sb, &queryParams, "@"+field, mapping.condition, mapping.value)
				}
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

func DeleteUserData(
	ctx context.Context,
	udContainer *azcosmos.ContainerClient,
	nsContainer *azcosmos.ContainerClient,
	userId string,
	city string,
) error {
	ctx1, cancel1 := context.WithTimeout(ctx, 10*time.Second)
	_, err := udContainer.DeleteItem(ctx1, azcosmos.NewPartitionKeyString(userId), strings.ToLower(city), nil)
	cancel1()
	if err != nil {
		return fmt.Errorf("failed to delete user data item: %w", err)
	}

	ctx2, cancel2 := context.WithTimeout(ctx, 10*time.Second)
	_, err = nsContainer.DeleteItem(ctx2, azcosmos.NewPartitionKeyString(strings.ToLower(city)), userId, nil)
	cancel2()
	if err != nil {
		return fmt.Errorf("failed to delete notification settings item: %w", err)
	}

	return nil
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
		"title":              {" AND (NOT IS_DEFINED(c.title) OR (CONTAINS(LOWER(@title), LOWER(c.title)) = true))", listing.Title},
		"housingCooperative": {" AND (NOT IS_DEFINED(c.housingCooperative) OR (CONTAINS(LOWER(@housingCooperative), LOWER(c.housingCooperative)) = true))", listing.HousingCooperative},
		"projectId":          {" AND (NOT IS_DEFINED(c.projectId) OR (LOWER(c.projectId) = LOWER(@projectId)))", listing.ProjectID},
		"postalCodes":        {" AND (NOT IS_DEFINED(c.postalCode) OR (CONTAINS(c.postalCode, @postalCodes) = true))", listing.PostalCode},
		"roomCount":          {" AND ((NOT IS_DEFINED(c.minRoomCount) OR (c.minRoomCount <= @roomCount)) AND (NOT IS_DEFINED(c.maxRoomCount) OR (c.maxRoomCount >= @roomCount)))", listing.RoomCount},
		"squareMeters":       {" AND ((NOT IS_DEFINED(c.minSqm) OR (c.minSqm <= @squareMeters)) AND (NOT IS_DEFINED(c.maxSqm) OR (c.maxSqm >= @squareMeters)))", listing.SquareMeters},
		"availabilityDate":   {" AND (NOT IS_DEFINED(c.availableFrom) OR (c.availableFrom <= @availabilityDate))", listing.AvailabilityDate},
		"yearBuilt":          {" AND ((NOT IS_DEFINED(c.minYearBuilt) OR (c.minYearBuilt <= @yearBuilt)) AND (NOT IS_DEFINED(c.maxYearBuilt) OR (c.maxYearBuilt >= @yearBuilt)))", listing.YearBuilt},
		"hwgEnergyClass":     {" AND (NOT IS_DEFINED(c.minHwgEnergyClass) OR (ARRAY_CONTAINS(@hwgEnergyClass, c.minHwgEnergyClass) = true))", listing.HwgEnergyClass},
		"fgeeEnergyClass":    {" AND (NOT IS_DEFINED(c.minFgeeEnergyClass) OR (ARRAY_CONTAINS(@fgeeEnergyClass, c.minFgeeEnergyClass) = true))", listing.FgeeEnergyClass},
		"listingType":        {" AND (NOT IS_DEFINED(c.listingType) OR c.listingType = 'both' OR @listingType = 'both' OR c.listingType = @listingType)", listing.ListingType},
		"rentPricePerMonth":  {" AND ((NOT IS_DEFINED(c.minRentPrice) OR (c.minRentPrice <= @rentPricePerMonth)) AND (NOT IS_DEFINED(c.maxRentPrice) OR (c.maxRentPrice >= @rentPricePerMonth)))", listing.RentPricePerMonth},
		"cooperativeShare":   {" AND ((NOT IS_DEFINED(c.minCooperativeShare) OR (c.minCooperativeShare <= @cooperativeShare)) AND (NOT IS_DEFINED(c.maxCooperativeShare) OR (c.maxCooperativeShare >= @cooperativeShare)))", listing.CooperativeShare},
		"salePrice":          {" AND ((NOT IS_DEFINED(c.minSalePrice) OR (c.minSalePrice <= @salePrice)) AND (NOT IS_DEFINED(c.maxSalePrice) OR (c.maxSalePrice >= @salePrice)))", listing.SalePrice},
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
				addQueryParam(&sb, &queryParams, "@"+field, mapping.condition, mapping.value)
			}
		case *string:
			if mapping.value != nil {
				ecStr, ok := (mapping.value).(*string)
				if !ok {
					return nil, fmt.Errorf("value of %s has incorrect format", field)
				}
				if ecStr != nil {
					if field == "hwgEnergyClass" || field == "fgeeEnergyClass" {
						ecClass := models.EnergyClass(*ecStr)
						addQueryParam(&sb, &queryParams, "@"+field, mapping.condition, models.GetEnergyClasses()[ecClass.GetIndex():])
					} else {
						addQueryParam(&sb, &queryParams, "@"+field, mapping.condition, *ecStr)
					}
				}
			}
		case int:
			addQueryParam(&sb, &queryParams, "@"+field, mapping.condition, mapping.value)
		case float32:
			v := mapping.value.(float32)
			addQueryParam(&sb, &queryParams, "@"+field, mapping.condition, v)
		case *int:
			if mapping.value != nil {
				v := mapping.value.(*int)
				if v != nil {
					addQueryParam(&sb, &queryParams, "@"+field, mapping.condition, v)
				}
			}
		case *float32:
			if mapping.value != nil {
				v := mapping.value.(*float32)
				if v != nil {
					addQueryParam(&sb, &queryParams, "@"+field, mapping.condition, *v)
				}
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
