package types

import "time"

type CatalogueItems struct {
	Items                []Item
	DominantBrand        interface{} `json:"dominant_brand"`
	SearchTrackingParams struct {
		SearchCorrelationID string `json:"search_correlation_id"`
		SearchSessionID     string `json:"search_session_id"`
	} `json:"search_tracking_params"`
	Pagination struct {
		CurrentPage  int `json:"current_page"`
		TotalPages   int `json:"total_pages"`
		TotalEntries int `json:"total_entries"`
		PerPage      int `json:"per_page"`
		Time         int `json:"time"`
	} `json:"pagination"`
	Code int `json:"code"`
}

type Item struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	Price string `json:"price"`

	IsVisible  int         `json:"is_visible"`
	Discount   interface{} `json:"discount"`
	BrandTitle string      `json:"brand_title"`
	User       struct {
		ID         int    `json:"id"`
		Login      string `json:"login"`
		Business   bool   `json:"business"`
		ProfileURL string `json:"profile_url"`
		Photo      struct {
			ID                  int         `json:"id"`
			Width               int         `json:"width"`
			Height              int         `json:"height"`
			TempUUID            interface{} `json:"temp_uuid"`
			URL                 string      `json:"url"`
			DominantColor       string      `json:"dominant_color"`
			DominantColorOpaque string      `json:"dominant_color_opaque"`
			Thumbnails          []struct {
				Type         string      `json:"type"`
				URL          string      `json:"url"`
				Width        int         `json:"width"`
				Height       int         `json:"height"`
				OriginalSize interface{} `json:"original_size"`
			} `json:"thumbnails"`
			IsSuspicious   bool        `json:"is_suspicious"`
			Orientation    interface{} `json:"orientation"`
			HighResolution struct {
				ID          string      `json:"id"`
				Timestamp   int         `json:"timestamp"`
				Orientation interface{} `json:"orientation"`
			} `json:"high_resolution"`
			FullSizeURL string `json:"full_size_url"`
			IsHidden    bool   `json:"is_hidden"`
			Extra       struct {
			} `json:"extra"`
		} `json:"photo"`
	} `json:"user"`
	URL      string `json:"url"`
	Promoted bool   `json:"promoted"`
	Photo    struct {
		ID                  int64  `json:"id"`
		ImageNo             int    `json:"image_no"`
		Width               int    `json:"width"`
		Height              int    `json:"height"`
		DominantColor       string `json:"dominant_color"`
		DominantColorOpaque string `json:"dominant_color_opaque"`
		URL                 string `json:"url"`
		IsMain              bool   `json:"is_main"`
		Thumbnails          []struct {
			Type         string      `json:"type"`
			URL          string      `json:"url"`
			Width        int         `json:"width"`
			Height       int         `json:"height"`
			OriginalSize interface{} `json:"original_size"`
		} `json:"thumbnails"`
		HighResolution struct {
			ID          string      `json:"id"`
			Timestamp   int         `json:"timestamp"`
			Orientation interface{} `json:"orientation"`
		} `json:"high_resolution"`
		IsSuspicious bool   `json:"is_suspicious"`
		FullSizeURL  string `json:"full_size_url"`
		IsHidden     bool   `json:"is_hidden"`
		Extra        struct {
		} `json:"extra"`
	} `json:"photo"`
	FavouriteCount int         `json:"favourite_count"`
	IsFavourite    bool        `json:"is_favourite"`
	Badge          interface{} `json:"badge"`
	Conversion     interface{} `json:"conversion"`

	TotalItemPrice        string        `json:"total_item_price"`
	TotalItemPriceRounded interface{}   `json:"total_item_price_rounded"`
	ViewCount             int           `json:"view_count"`
	SizeTitle             string        `json:"size_title"`
	ContentSource         string        `json:"content_source"`
	Status                string        `json:"status"`
	IconBadges            []interface{} `json:"icon_badges"`
	SearchTrackingParams  struct {
		Score          float64  `json:"score"`
		MatchedQueries []string `json:"matched_queries"`
	} `json:"search_tracking_params"`
	Timestamp  time.Time
	StringTime string
	TimeDiff   time.Duration
}

type ItemDetails struct {
	ID                     int       `json:"id"`
	Title                  string    `json:"title"`
	BrandID                int       `json:"brand_id"`
	SizeID                 int       `json:"size_id"`
	StatusID               int       `json:"status_id"`
	UserID                 int       `json:"user_id"`
	CountryID              int       `json:"country_id"`
	CatalogID              int       `json:"catalog_id"`
	Color1ID               int       `json:"color1_id"`
	Color2ID               *int      `json:"color2_id"` // Pointer type for nullable field
	PackageSizeID          int       `json:"package_size_id"`
	IsVisible              bool      `json:"is_visible"`
	IsUnisex               bool      `json:"is_unisex"`
	IsClosed               bool      `json:"is_closed"`
	ModerationStatus       int       `json:"moderation_status"`
	IsHidden               bool      `json:"is_hidden"`
	FavouriteCount         int       `json:"favourite_count"`
	ActiveBidCount         int       `json:"active_bid_count"`
	Description            string    `json:"description"`
	PackageSizeStandard    bool      `json:"package_size_standard"`
	ItemClosingAction      *string   `json:"item_closing_action"` // Pointer type for nullable field
	RelatedCatalogIDs      []int     `json:"related_catalog_ids"`
	RelatedCatalogsEnabled bool      `json:"related_catalogs_enabled"`
	Size                   string    `json:"size"`
	Brand                  string    `json:"brand"`
	Composition            string    `json:"composition"`
	ExtraConditions        string    `json:"extra_conditions"`
	DisposalConditions     int       `json:"disposal_conditions"`
	IsForSell              bool      `json:"is_for_sell"`
	IsHandicraft           bool      `json:"is_handicraft"`
	IsProcessing           bool      `json:"is_processing"`
	IsDraft                bool      `json:"is_draft"`
	IsReserved             bool      `json:"is_reserved"`
	Label                  string    `json:"label"`
	OriginalPriceNumeric   string    `json:"original_price_numeric"`
	Currency               string    `json:"currency"`
	PriceNumeric           string    `json:"price_numeric"`
	LastPushUpAt           time.Time `json:"last_push_up_at"`
	LastPushUpAtNew        time.Time `json:"last_push_up_at_new"`
	CreatedAtTS            time.Time `json:"created_at_ts"`
	UpdatedAtTS            time.Time `json:"updated_at_ts"`
	UserUpdatedAtTS        time.Time `json:"user_updated_at_ts"`
	IsDelayedPublication   bool      `json:"is_delayed_publication"`

	CanBeSold               bool       `json:"can_be_sold"`
	CanFeedback             bool       `json:"can_feedback"`
	ItemReservationID       *int       `json:"item_reservation_id"`      // Pointer type for nullable field
	PromotedUntil           *time.Time `json:"promoted_until"`           // Pointer type for nullable field
	PromotedInternationally *bool      `json:"promoted_internationally"` // Pointer type for nullable field
	DiscountPriceNumeric    *string    `json:"discount_price_numeric"`   // Pointer type for nullable field
	Author                  *string    `json:"author"`                   // Pointer type for nullable field
	BookTitle               *string    `json:"book_title"`               // Pointer type for nullable field
	ISBN                    *string    `json:"isbn"`                     // Pointer type for nullable field
	MeasurementWidth        *string    `json:"measurement_width"`        // Pointer type for nullable field
	MeasurementLength       *string    `json:"measurement_length"`       // Pointer type for nullable field
	MeasurementUnit         *string    `json:"measurement_unit"`         // Pointer type for nullable field
	Manufacturer            *string    `json:"manufacturer"`             // Pointer type for nullable field
	ManufacturerLabelling   *string    `json:"manufacturer_labelling"`   // Pointer type for nullable field
	TransactionPermitted    bool       `json:"transaction_permitted"`
	VideoGameRatingID       *int       `json:"video_game_rating_id"` // Pointer type for nullable field
	ItemAttributes          []string   `json:"item_attributes"`
	HaovItem                bool       `json:"haov_item?"`
}
