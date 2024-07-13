package types

import "time"

type Subscription struct {
	ID          int
	Url         string //this will become the preferences
	Webhook     string
	Preferences map[string]interface{}
}
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
type Region struct {
	BaseUrl  string
	Currency string
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
	ID                      int                   `json:"id"`
	Title                   string                `json:"title"`
	BrandID                 int                   `json:"brand_id"`
	SizeID                  int                   `json:"size_id"`
	StatusID                int                   `json:"status_id"`
	UserID                  int                   `json:"user_id"`
	CountryID               int                   `json:"country_id"`
	CatalogID               int                   `json:"catalog_id"`
	Color1ID                int                   `json:"color1_id"`
	Color2ID                *int                  `json:"color2_id"`
	PackageSizeID           int                   `json:"package_size_id"`
	IsVisible               int                   `json:"is_visible"`
	IsUnisex                int                   `json:"is_unisex"`
	ModerationStatus        int                   `json:"moderation_status"`
	IsHidden                bool                  `json:"is_hidden"`
	IsClosed                int                   `json:"is_closed"`
	IsClosedNew             bool                  `json:"is_closed_new"`
	FavouriteCount          int                   `json:"favourite_count"`
	ActiveBidCount          int                   `json:"active_bid_count"`
	Description             string                `json:"description"`
	PackageSizeStandard     bool                  `json:"package_size_standard"`
	ItemClosingAction       interface{}           `json:"item_closing_action"`
	RelatedCatalogIDs       []interface{}         `json:"related_catalog_ids"`
	RelatedCatalogsEnabled  bool                  `json:"related_catalogs_enabled"`
	Size                    string                `json:"size"`
	Brand                   string                `json:"brand"`
	Composition             string                `json:"composition"`
	ExtraConditions         string                `json:"extra_conditions"`
	DisposalConditions      int                   `json:"disposal_conditions"`
	IsForSell               bool                  `json:"is_for_sell"`
	IsHandicraft            bool                  `json:"is_handicraft"`
	IsProcessing            bool                  `json:"is_processing"`
	IsDraft                 bool                  `json:"is_draft"`
	IsReserved              bool                  `json:"is_reserved"`
	Label                   string                `json:"label"`
	OriginalPriceNumeric    string                `json:"original_price_numeric"`
	Currency                string                `json:"currency"`
	PriceNumeric            string                `json:"price_numeric"`
	LastPushUpAt            string                `json:"last_push_up_at"`
	CreatedAtTs             string                `json:"created_at_ts"`
	UpdatedAtTs             string                `json:"updated_at_ts"`
	UserUpdatedAtTs         string                `json:"user_updated_at_ts"`
	IsDelayedPublication    bool                  `json:"is_delayed_publication"`
	Photos                  []Photo               `json:"photos"`
	CanBeSold               bool                  `json:"can_be_sold"`
	CanFeedback             bool                  `json:"can_feedback"`
	PossibleToCloseNew      bool                  `json:"possible_to_close_new"`
	ItemReservationID       interface{}           `json:"item_reservation_id"`
	PromotedUntil           interface{}           `json:"promoted_until"`
	PromotedInternationally interface{}           `json:"promoted_internationally"`
	DiscountPriceNumeric    interface{}           `json:"discount_price_numeric"`
	Author                  interface{}           `json:"author"`
	BookTitle               interface{}           `json:"book_title"`
	ISBN                    interface{}           `json:"isbn"`
	MeasurementWidth        interface{}           `json:"measurement_width"`
	MeasurementLength       interface{}           `json:"measurement_length"`
	MeasurementUnit         interface{}           `json:"measurement_unit"`
	Manufacturer            interface{}           `json:"manufacturer"`
	ManufacturerLabelling   interface{}           `json:"manufacturer_labelling"`
	TransactionPermitted    bool                  `json:"transaction_permitted"`
	VideoGameRatingID       interface{}           `json:"video_game_rating_id"`
	ItemAttributes          []interface{}         `json:"item_attributes"`
	HaovItem                bool                  `json:"haov_item"`
	User                    struct{}              `json:"user"`
	Price                   Price                 `json:"price"`
	DiscountPrice           interface{}           `json:"discount_price"`
	ServiceFee              string                `json:"service_fee"`
	TotalItemPrice          string                `json:"total_item_price"`
	CanEdit                 bool                  `json:"can_edit"`
	CanDelete               bool                  `json:"can_delete"`
	CanReserve              bool                  `json:"can_reserve"`
	CanMarkAsSold           bool                  `json:"can_mark_as_sold"`
	CanTransfer             bool                  `json:"can_transfer"`
	InstantBuy              bool                  `json:"instant_buy"`
	CanClose                bool                  `json:"can_close"`
	CanBuy                  bool                  `json:"can_buy"`
	CanBundle               bool                  `json:"can_bundle"`
	CanAskSeller            bool                  `json:"can_ask_seller"`
	CanFavourite            bool                  `json:"can_favourite"`
	UserLogin               string                `json:"user_login"`
	CityID                  interface{}           `json:"city_id"`
	City                    string                `json:"city"`
	Country                 string                `json:"country"`
	Promoted                bool                  `json:"promoted"`
	IsMobile                bool                  `json:"is_mobile"`
	BumpBadgeVisible        bool                  `json:"bump_badge_visible"`
	BrandDto                BrandDto              `json:"brand_dto"`
	CatalogBranchTitle      string                `json:"catalog_branch_title"`
	Path                    string                `json:"path"`
	URL                     string                `json:"url"`
	AcceptedPayInMethods    []AcceptedPayInMethod `json:"accepted_pay_in_methods"`
	CreatedAt               string                `json:"created_at"`
	Color1                  string                `json:"color1"`
	Color2                  interface{}           `json:"color2"`
	SizeTitle               string                `json:"size_title"`
	DescriptionAttributes   []interface{}         `json:"description_attributes"`
	VideoGameRating         interface{}           `json:"video_game_rating"`
	Status                  string                `json:"status"`
	IsFavourite             bool                  `json:"is_favourite"`
	ViewCount               int                   `json:"view_count"`
	Performance             interface{}           `json:"performance"`
	StatsVisible            bool                  `json:"stats_visible"`
	CanPushUp               bool                  `json:"can_push_up"`
	Badge                   interface{}           `json:"badge"`
	SizeGuideFAQEntryID     int                   `json:"size_guide_faq_entry_id"`
	Localization            string                `json:"localization"`
	OfflineVerification     bool                  `json:"offline_verification"`
	OfflineVerificationFee  interface{}           `json:"offline_verification_fee"`
	IconBadges              []interface{}         `json:"icon_badges"`
	StringTime              string
}

type Photo struct {
	ID                  int            `json:"id"`
	ImageNo             int            `json:"image_no"`
	Width               int            `json:"width"`
	Height              int            `json:"height"`
	DominantColor       string         `json:"dominant_color"`
	DominantColorOpaque string         `json:"dominant_color_opaque"`
	URL                 string         `json:"url"`
	IsMain              bool           `json:"is_main"`
	Thumbnails          []Thumbnail    `json:"thumbnails"`
	HighResolution      HighResolution `json:"high_resolution"`
	IsSuspicious        bool           `json:"is_suspicious"`
	FullSizeURL         string         `json:"full_size_url"`
	IsHidden            bool           `json:"is_hidden"`
	Extra               struct{}       `json:"extra"`
}

type Thumbnail struct {
	Type         string      `json:"type"`
	URL          string      `json:"url"`
	Width        int         `json:"width"`
	Height       int         `json:"height"`
	OriginalSize interface{} `json:"original_size"`
}

type HighResolution struct {
	ID          string      `json:"id"`
	Timestamp   int         `json:"timestamp"`
	Orientation interface{} `json:"orientation"`
}

type Price struct {
	Amount       string `json:"amount"`
	CurrencyCode string `json:"currency_code"`
}

type BrandDto struct {
	ID                        int    `json:"id"`
	Title                     string `json:"title"`
	Slug                      string `json:"slug"`
	FavouriteCount            int    `json:"favourite_count"`
	PrettyFavouriteCount      string `json:"pretty_favourite_count"`
	ItemCount                 int    `json:"item_count"`
	PrettyItemCount           string `json:"pretty_item_count"`
	IsVisibleInListings       bool   `json:"is_visible_in_listings"`
	RequiresAuthenticityCheck bool   `json:"requires_authenticity_check"`
	IsLuxury                  bool   `json:"is_luxury"`
	IsHVF                     bool   `json:"is_hvf"`
	Path                      string `json:"path"`
	URL                       string `json:"url"`
	IsFavourite               bool   `json:"is_favourite"`
}

type AcceptedPayInMethod struct {
	ID                   int    `json:"id"`
	Code                 string `json:"code"`
	RequiresCreditCard   bool   `json:"requires_credit_card"`
	EventTrackingCode    string `json:"event_tracking_code"`
	Icon                 string `json:"icon"`
	Enabled              bool   `json:"enabled"`
	TranslatedName       string `json:"translated_name"`
	Note                 string `json:"note"`
	MethodChangePossible bool   `json:"method_change_possible"`
}

type Response struct {
	Item ItemDetails `json:"item"`
	Code int         `json:"code"`
}
