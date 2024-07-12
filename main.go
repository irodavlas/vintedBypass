package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/vintedMonitor/utils"
)

const (
	UK_BASE_URL = "https://www.vinted.co.uk"
	FR_BASE_URL = "https://www.vinted.fr"
	DE_BASE_URL = "https://www.vinted.de"
	PL_BASE_URL = "https://www.vinted.pl"
)
const UK_API_URL = UK_BASE_URL + "/api/v2/items/"

func main() {
	url := "https://www.vinted.co.uk/catalog?catalog[]=84&brand_ids[]=88&status_ids[]=2&status_ids[]=3&price_to=15&currency=GBP&color_ids[]=27&color_ids[]=9&order=newest_first"
	filterData, err := utils.ParseURLParameters(url)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return
	}

	filterDict := utils.CreateFilterDict(filterData)
	for key, value := range filterDict {
		fmt.Printf("%s: %v\n", key, value)
	}
	/*
		monitor := utils.Latest_Sku_Monitor{
			Latest_channel: make(chan int, 99999),
			Proxies:        parse_proxy_file(),
		}
		go func() {
			for {
				select {
				case sku := <-monitor.Latest_channel:
					log.Println(sku)

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
		select {}
	*/
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
