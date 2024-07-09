package main

import (
	"log"
	"sync"

	utils "github.com/vintedMonitor/utils"
)

const (
	UK_BASE_URL = "https://www.vinted.co.uk"
	FR_BASE_URL = "https://www.vinted.fr"
	DE_BASE_URL = "https://www.vinted.de"
	PL_BASE_URL = "https://www.vinted.pl"
)
const UK_API_URL = UK_BASE_URL + "/api/v2/items/"

func main() {

	monitor := utils.Latest_Sku_Monitor{
		Latest_channel: make(chan string),
		Latest_sku:     "",
		LatestMux:      sync.Mutex{},
	}
	go func() {
		for {
			select {
			case latest_sku := <-monitor.Latest_channel:
				monitor.LatestMux.Lock()
				monitor.Latest_sku = latest_sku
				monitor.LatestMux.Unlock()
			case <-monitor.New_batch_signal:
				//new batch start
				monitor.Start_new_batch()

			}

		}
	}()
	client, err := utils.NewClient("https://www.vinted.com")
	if err != nil {
		log.Fatal(err)
	}
	//start mointor for starting page
	monitor.Session, err = monitor.Get_session_cookie(client)
	if err != nil {
		log.Fatal(err)
	}
	monitor.Get_latest_sku(client, monitor.Session, monitor.Latest_channel)

	select {}

}
