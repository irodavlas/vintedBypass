package utils

import (
	"net/url"
	"strconv"
)

type FilterData struct {
	BrandIDs    []int
	Catalog     []int
	PriceFrom   *int
	PriceTo     *int
	Currency    *string
	ColorIDs    []int
	SizeIDs     []int
	MaterialIDs []int
	StatusIDs   []int
	SearchText  *string
}

func ParseURLParameters(data string) (*FilterData, error) {
	parsedURL, err := url.Parse(data)
	if err != nil {
		return nil, err
	}

	params := parsedURL.Query()

	parseIntArray := func(key string) []int {
		values := params[key]
		result := make([]int, 0, len(values))
		for _, v := range values {
			if intValue, err := strconv.Atoi(v); err == nil {
				result = append(result, intValue)
			}
		}
		return result
	}

	parseInt := func(key string) *int {
		if value := params.Get(key); value != "" {
			if intValue, err := strconv.Atoi(value); err == nil {
				return &intValue
			}
		}
		return nil
	}

	parseString := func(key string) *string {
		if value := params.Get(key); value != "" {
			return &value
		}
		return nil
	}

	filterData := &FilterData{
		BrandIDs:    parseIntArray("brand_ids[]"),
		Catalog:     parseIntArray("catalog[]"),
		PriceFrom:   parseInt("price_from"),
		PriceTo:     parseInt("price_to"),
		Currency:    parseString("currency"),
		ColorIDs:    parseIntArray("color_ids[]"),
		SizeIDs:     parseIntArray("size_ids[]"),
		MaterialIDs: parseIntArray("material_ids[]"),
		StatusIDs:   parseIntArray("status_ids[]"),
		SearchText:  parseString("search_text"),
	}

	return filterData, nil
}

func CreateFilterDict(filter *FilterData) map[string]interface{} {
	filterDict := make(map[string]interface{})

	setIfNotNull := func(key string, value interface{}, name string) {
		if value != nil {
			switch v := value.(type) {
			case []int:
				if len(v) > 0 {
					filterDict[name] = v
				}
			case *string:
				if *v != "" {
					filterDict[name] = *v
				}
			case *int:
				filterDict[name] = *v
			}
		}
	}

	setIfNotNull("color_ids", filter.ColorIDs, "color_ids")
	setIfNotNull("brand_ids", filter.BrandIDs, "brand_ids")
	setIfNotNull("size_ids", filter.SizeIDs, "size_ids")
	setIfNotNull("material_ids", filter.MaterialIDs, "material_ids")
	setIfNotNull("status_ids", filter.StatusIDs, "status_ids")
	setIfNotNull("catalog", filter.Catalog, "catalog")
	setIfNotNull("search_text", filter.SearchText, "search_text")
	setIfNotNull("currency", filter.Currency, "currency")

	return filterDict
}
