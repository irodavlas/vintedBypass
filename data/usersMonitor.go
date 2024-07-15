package data

import (
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/vintedMonitor/database"
	"github.com/vintedMonitor/types"
	"github.com/vintedMonitor/utils"
	"github.com/vintedMonitor/webhook"
)

type Monitor struct {
	Id              int
	Username        string
	Item_channel    chan types.ItemDetails
	Subscriptions   []types.Subscription
	SubscriptionMux sync.Mutex
}

const (
	UK_BASE_URL = "https://www.vinted.co.uk"
	FR_BASE_URL = "https://www.vinted.fr"
	DE_BASE_URL = "https://www.vinted.de"
	PL_BASE_URL = "https://www.vinted.pl"
)
const (
	UK_CURRENCY = "£"
	FR_CURRENCY = "€"
	DE_CURRENCY = "€"
	PL_CURRENCY = "zł"
)

var Regions = map[int]types.Region{
	16: {BaseUrl: FR_BASE_URL, Currency: FR_CURRENCY},
	13: {BaseUrl: UK_BASE_URL, Currency: UK_CURRENCY},
	2:  {BaseUrl: DE_BASE_URL, Currency: DE_CURRENCY},
	15: {BaseUrl: PL_BASE_URL, Currency: PL_CURRENCY},
}

func (m *Monitor) Start_user_dispatcher(db *database.MyDB) {
	log.Println("Starting dispatcher for user:", m.Username)
	time.Sleep(2 * time.Second)
	go func() {
		{
			ticker := time.NewTicker(180 * time.Second) // Adjust interval as needed
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					subs, err := m.fetchSubscriptions(db)
					if err != nil {
						log.Println("Error occurred while fetching user Preferences: ", err)
					}
					m.SubscriptionMux.Lock()
					m.Subscriptions = subs
					m.SubscriptionMux.Unlock()
					log.Println("Succesfully updated preferences for user:", m.Username)
				}

			}
		}
	}()
	for {
		select {
		case item := <-m.Item_channel:

			go m.check_keywords(item)

		}
	}
}
func Create_user_dispatcher(id int, username string, db *database.MyDB) (*Monitor, error) {
	monitor := Monitor{Id: id, Username: username, Item_channel: make(chan types.ItemDetails, 9999)}
	subs, err := monitor.fetchSubscriptions(db)
	if err != nil {
		return nil, err
	}
	monitor.Subscriptions = subs

	//start go routine to listen to pid and add keywords
	return &monitor, nil
}

// fetch subs
// update subs (remove subs)

func (m *Monitor) check_keywords(item types.ItemDetails) {
	for _, sub := range m.Subscriptions {
		if region, exists := Regions[item.CountryID]; exists {
			if matchesFilter(item, sub.Preferences) {
				//send webhook
				log.Printf("New Item monitored on [%d], PID:[%d]", m.Id, item.ID)
				webhook.Send_webhook(sub.Webhook, region, item)
				break
			}
		}

	}
}
func matchesFilter(item types.ItemDetails, filters map[string]interface{}) bool {
	if brandIDs, ok := filters["brand_ids"].([]int); ok && brandIDs != nil && !containsInt(brandIDs, item.BrandID) {
		return false
	}

	if catalog, ok := filters["catalog"].([]int); ok && catalog != nil && !containsInt(catalog, item.CatalogID) {
		return false
	}

	if priceFrom, ok := filters["price_from"].(int); ok && priceFrom != 0 {
		price := parsePrice(item.PriceNumeric)
		if price < priceFrom {
			return false
		}
	}

	if priceTo, ok := filters["price_to"].(int); ok && priceTo != 0 {
		price := parsePrice(item.PriceNumeric)
		if price > priceTo {
			return false
		}
	}

	if currency, ok := filters["currency"].(string); ok && currency != "" && item.Currency != currency {
		return false
	}

	if statusIDs, ok := filters["status_ids"].([]int); ok && statusIDs != nil && !containsInt(statusIDs, item.StatusID) {
		return false
	}

	if searchText, ok := filters["search_text"].(string); ok && searchText != "" && !strings.Contains(strings.ToLower(item.UserLogin), strings.ToLower(searchText)) {
		return false
	}

	return true
}

// Helper function to check if a slice contains a given integer
func containsInt(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// Helper function to parse price from string to int
func parsePrice(priceStr string) int {
	price, err := strconv.Atoi(priceStr)
	if err != nil {
		return 0
	}
	return price
}

// each tot fewtch subs if none are found stop monitor and remove it from map
func (m *Monitor) fetchSubscriptions(db *database.MyDB) ([]types.Subscription, error) {
	subs, err := db.Get_Subscription_of_user(m.Username)
	if err != nil {
		return nil, err
	}
	var final_subs []types.Subscription
	m.SubscriptionMux.Lock()
	for _, sub := range subs {
		sub.Preferences = utils.Filter_user_subscription(sub.Url)
		final_subs = append(final_subs, sub)
	}
	defer m.SubscriptionMux.Unlock()
	return final_subs, nil

}
func (m *Monitor) update_subscriptions(db *database.MyDB) {

}
