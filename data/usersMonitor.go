package data

import (
	"sync"

	"github.com/vintedMonitor/database"
	"github.com/vintedMonitor/types"
	"github.com/vintedMonitor/utils"
)

type Monitor struct {
	Id       int
	Username string

	Subscriptions   []types.Subscription
	SubscriptionMux sync.Mutex
}

func Start_user_dispatcher(id int, username string, db *database.MyDB) (*Monitor, error) {
	monitor := Monitor{Id: id, Username: username}
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
