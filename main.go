package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
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
const BATCH_SIZE = 1000

func main() {

	monitor := utils.Latest_Sku_Monitor{
		Latest_channel:   make(chan utils.Sku, 99999),
		New_batch_signal: make(chan bool, 100),
		Latest_sku:       0,
		LatestMux:        sync.Mutex{},
		Proxies:          parse_proxy_file(),
	}
	go func() {
		for {
			select {
			case item := <-monitor.Latest_channel:
				//log.Printf("sku:%d Time:%s", item.Sku, item.Time)
				monitor.LatestMux.Lock()
				//put here logic to start new batches
				//use another value such as last round sku
				if item.Sku > monitor.Latest_sku {
					monitor.Latest_sku = item.Sku
				}

				monitor.Latest_batch_sku = monitor.Latest_sku
				monitor.LatestMux.Unlock()

				//go monitor.Start_new_batch(monitor.Latest_batch_sku, 100)
			case <-monitor.New_batch_signal:
				go monitor.Start_new_batch(monitor.Latest_sku, BATCH_SIZE)
			}
			//case start new batch once all terminate

		}
	}()
	go func() {
		for {
			select {}
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
	monitor.Latest_batch_sku = monitor.Get_latest_sku(client, monitor.Session)
	print(monitor.Latest_batch_sku)
	monitor.Start_new_batch(monitor.Latest_batch_sku, BATCH_SIZE)
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
