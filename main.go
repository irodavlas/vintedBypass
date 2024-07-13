package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/vintedMonitor/data"
	"github.com/vintedMonitor/database"
	"github.com/vintedMonitor/types"
)

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

func main() {
	monitors := make(map[string]*data.Monitor)
	db, err := database.Connect_to_database()
	if err != nil {
		db.DB.Close()
		log.Fatal("Error connecting to DB: ", err)
	}
	users, err := db.Get_all_users()
	if err != nil {
		log.Fatal("Error retrieving users:", err)
	}
	for i, user := range *users {
		monitor, err := data.Start_user_dispatcher(i, user, db)
		if err != nil {
			log.Fatalf("Error creating monitor for user: %s, with error: %s", user, err)
		}
		monitors[user] = monitor
	}
	/*
		monitor := utils.Latest_Sku_Monitor{
			Latest_channel: make(chan types.ItemDetails, 99999),
			Proxies:        parse_proxy_file(),
		}
		go func() {
			for {
				select {
				case item := <-monitor.Latest_channel:
					log.Println(item.Country, item.CountryID)

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
		// need a goroutine that updated the session once in a while

		monitor.Start_monitor()
	*/
	select {}

}
func formatProxy(proxy string) string {
	// Split the proxy string by colon
	parts := strings.Split(proxy, ":")

	// Check if the proxy has the correct number of parts
	if len(parts) != 4 {
		return "Invalid proxy format"
	}

	// Extract host, port, user, and pass
	host := parts[0]
	port := parts[1]
	user := parts[2]
	pass := parts[3]

	// Format the proxy in the required format
	formattedProxy := fmt.Sprintf("http://%s:%s@%s:%s", user, pass, host, port)
	return formattedProxy
}

func parse_proxy_file() []string {
	var proxies []string
	file, err := os.Open("proxy.txt")
	if err != nil {
		log.Fatalf("Error loading proxies file: %s", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		proxy := formatProxy(line)
		proxies = append(proxies, proxy)
	}
	return proxies
}
