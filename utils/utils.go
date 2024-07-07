package utils

import (
    "fmt"
    "net/url"
    "strconv"
)

// Struct to hold the parsed parameters
type FilterParams struct {
    BrandIds    []int
    Catalog     []int
    PriceFrom   int
    PriceTo     int
    ColorIds    []int
    SizeIds     []int
    MaterialIds []int
    StatusIds   []int
    StatusId    []int
    CatalogId   []int
    BrandId     []int
    SearchText  string
    Currency    string
}

// Function to parse URL parameters
func parseURLParameters(rawURL string) (FilterParams, error) {
    u, err := url.Parse(rawURL)
    if err != nil {
        return FilterParams{}, err
    }

    params := u.Query()

    parseIntSlice := func(values []string) []int {
        var result []int
        for _, v := range values {
            if iv, err := strconv.Atoi(v); err == nil {
                result = append(result, iv)
            }
        }
        return result
    }

    parseInt := func(value string) int {
        if iv, err := strconv.Atoi(value); err == nil {
            return iv
        }
        return 0
    }

    return FilterParams{
        BrandIds:    parseIntSlice(params["brand_ids[]"]),
        Catalog:     parseIntSlice(params["catalog[]"]),
        PriceFrom:   parseInt(params.Get("price_from")),
        PriceTo:     parseInt(params.Get("price_to")),
        ColorIds:    parseIntSlice(params["color_ids[]"]),
        SizeIds:     parseIntSlice(params["size_ids[]"]),
        MaterialIds: parseIntSlice(params["material_ids[]"]),
        StatusIds:   parseIntSlice(params["status_ids[]"]),
        StatusId:    parseIntSlice(params["status_id[]"]),
        CatalogId:   parseIntSlice(params["catalog_id[]"]),
        BrandId:     parseIntSlice(params["brand_id[]"]),
        SearchText:  params.Get("search_text"),
        Currency:    params.Get("currency"),
    }, nil
}

// Function to create filter map
func createFilterDict(filter FilterParams) map[string]interface{} {
    filterDict := make(map[string]interface{})

    setIfNotEmpty := func(key string, value interface{}) {
        switch v := value.(type) {
        case []int:
            if len(v) > 0 {
                filterDict[key] = v
            }
        case string:
            if v != "" {
                filterDict[key] = v
            }
        case int:
            if v != 0 {
                filterDict[key] = v
            }
        }
    }

    setIfNotEmpty("color_ids", filter.ColorIds)
    setIfNotEmpty("brand_ids", filter.BrandIds)
    setIfNotEmpty("size_ids", filter.SizeIds)
    setIfNotEmpty("material_ids", filter.MaterialIds)
    setIfNotEmpty("status_ids", filter.StatusIds)
    setIfNotEmpty("catalog", filter.Catalog)
    setIfNotEmpty("search_text", filter.SearchText)
    setIfNotEmpty("currency", filter.Currency)

    return filterDict
}

func main() {
    // Example usage
    rawURL := "https://www.vinted.co.uk/catalog?status_ids[]=2&currency=GBP&order=newest_first&price_from=8&catalog[]=1823&price_to=15&brand_ids[]=362&color_ids[]=3&color_ids[]=27"

    filterParams, err := parseURLParameters(rawURL)
    if err != nil {
        fmt.Println("Error parsing URL:", err)
        return
    }

    filterDict := createFilterDict(filterParams)

    fmt.Println("Filter Params:", filterParams)
    fmt.Println("Filter Dict:", filterDict)
}

