package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/vintedMonitor/data"
	"github.com/vintedMonitor/database"
	"github.com/vintedMonitor/types"
	"github.com/vintedMonitor/utils"
	"golang.org/x/exp/rand"
)

type User_monitors struct {
	Monitors   map[string]*data.Monitor
	MonitorMux sync.Mutex
}
type Users struct {
	Users []string
}

func main() {
	monitors := User_monitors{}
	monitors.Monitors = make(map[string]*data.Monitor)
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
		monitor, err := data.Create_user_dispatcher(i, user, db)
		if err != nil {
			log.Fatalf("Error creating monitor for user: %s, with error: %s", user, err)
		}
		monitors.Monitors[user] = monitor
	}

	for _, monitor := range monitors.Monitors {
		go monitor.Start_user_dispatcher(db)
	}

	scraping_monitor := utils.Latest_Sku_Monitor{
		Latest_channel: make(chan types.ItemDetails, 99999),
		Proxies:        parse_proxy_file(),
	}
	go func() {
		for {
			select {
			case item := <-scraping_monitor.Latest_channel:
				log.Println("New item found: ", item.ID)
				for _, monitor := range monitors.Monitors {
					monitor.Item_channel <- item
				}

			}

		}
	}()

	client, err := utils.NewClient("https://www.vinted.com")
	if err != nil {
		log.Fatal(err)
	}
	//start mointor for starting page
	scraping_monitor.Session, err = scraping_monitor.Get_session_cookie(client)
	if err != nil {
		log.Fatal(err)
	}
	// need a goroutine that updated the session once in a while

	scraping_monitor.Start_monitor()
	oldListOfUsers := Users{
		Users: *users,
	}

	//registers new user to the process, then i can delete the monitor elsewhere
	go func() {
		ticker := time.NewTicker(60 * time.Second) // Adjust interval as needed
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				newUsers, err := db.Get_all_users()
				if err != nil {
					log.Printf("Error retrieving users: %s", err)
					continue
				}
				missingUsers := findMissing(oldListOfUsers.Users, *newUsers)
				oldListOfUsers.Users = *newUsers
				for _, user := range missingUsers {
					log.Println("\n\n\nNew monitor created for user:", user)
					monitor, err := data.Create_user_dispatcher(rand.Intn(99999), user, db)
					if err != nil {
						log.Fatalf("Error creating monitor for user: %s, with error: %s", user, err)
					}
					monitors.Monitors[user] = monitor
					go monitor.Start_user_dispatcher(db)
				}

			}
		}
	}()
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

func findMissing(mainList, subList []string) []string {
	// Create a map to track occurrences in the main list
	mainMap := make(map[string]bool)

	// Populate the map with the main list values
	for _, item := range mainList {
		mainMap[item] = true
	}

	// Find elements in the sublist that are missing in the main list
	var missing []string
	for _, item := range subList {
		if !mainMap[item] {
			missing = append(missing, item)
		}
	}

	return missing
}
