package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
	"github.com/vintedMonitor/types"
	"golang.org/x/exp/rand"
)

type Latest_Sku_Monitor struct {
	Proxies        []string
	Latest_channel chan types.ItemDetails
	Session        string
}
type Latest_sku struct {
	Latest_sku int64
	//LatestMux  sync.Mutex
}
type Client struct {
	TlsClient *tls_client.HttpClient
	url       string
}

type Options struct {
	settings []tls_client.HttpClientOption
}

var MAX_RETRY = 200

func (m *Latest_Sku_Monitor) Get_session_cookie(session_client *Client) (string, error) {

	req, err := http.NewRequest(http.MethodGet, session_client.url, nil)
	if err != nil {
		return "", err

	}

	req.Header = http.Header{}
	req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Add("accept-language", "en-US,en;q=0.9")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("pragma", "no-cache")
	req.Header.Add("priority", "u=0, i")
	req.Header.Add("sec-ch-ua", "\"Google Chrome\";v=\"125\", \"Chromium\";v=\"125\", \"Not.A/Brand\";v=\"24\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "\"Windows\"")
	req.Header.Add("sec-fetch-dest", "document")
	req.Header.Add("sec-fetch-mode", "navigate")
	req.Header.Add("sec-fetch-site", "none")
	req.Header.Add("sec-fetch-user", "?1")
	req.Header.Add("upgrade-insecure-requests", "1")
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36")

	resp, err := (*session_client.TlsClient).Do(req)
	if err != nil {
		return "", err
	}

	resp.Body.Close()

	cookie := resp.Header["Set-Cookie"]
	session := extractSessionCookie(cookie)
	if session == "" {
		return "", err
	}

	log.Printf("---- Fetched Session cookie----")

	return session, nil

}
func extractSessionCookie(cookies []string) string {
	for _, cookie := range cookies {
		if strings.Contains(cookie, "_vinted_fr_session") {
			return strings.Split(strings.Split(cookie, "_vinted_fr_session=")[1], ";")[0]
		}
	}
	return ""
}
func (m *Latest_Sku_Monitor) Get_latest_sku(client *Client, session string) int64 {
	for {

		url := "https://www.vinted.com/api/v2/catalog/items?page=1&per_page=96&search_text=&catalog_ids=&order=newest_first&size_ids=&brand_ids=&status_ids=&color_ids=&material_ids="
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			log.Println("Retrying, Error occured: ", err)

			continue
		}

		req.Header = http.Header{}

		req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
		req.Header.Add("accept-language", "en-US,en;q=0.9")
		req.Header.Add("cookie", fmt.Sprintf("_vinted_fr_session=%s", session))
		req.Header.Add("sec-ch-ua", "\"Google Chrome\";v=\"125\", \"Chromium\";v=\"125\", \"Not.A/Brand\";v=\"24\"")
		req.Header.Add("sec-ch-ua-mobile", "?0")
		req.Header.Add("sec-ch-ua-platform", "\"Windows\"")
		req.Header.Add("upgrade-insecure-requests", "1")
		req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36")

		resp, err := (*client.TlsClient).Do(req)
		if err != nil {

			log.Println("Retrying, Error occured: ", err)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Retrying, Error occured: ", err)
			resp.Body.Close()
			continue
		}
		defer resp.Body.Close()

		log.Printf("(%d) - %s", resp.StatusCode, url)

		var data types.CatalogueItems
		err = json.Unmarshal(body, &data)
		if err != nil {
			log.Println("Error unmarshalling into CatalogueItems: ", err)
			continue
		}

		return data.Items[0].ID
	}

}

// makes the item requests in a for loop
func Make_request(sku int, client Client, monitor *Latest_Sku_Monitor, global_pid_list *Latest_sku) {

	var last_pid = sku
	var session = monitor.Session

	var retry_count = 0
	for {

		retry_count++
		if retry_count >= MAX_RETRY {
			log.Println("Dropping sku: ", last_pid)
			last_pid = int(global_pid_list.get_new_pid())
			retry_count = 0
			continue
		}
		url := "https://www.vinted.com/api/v2/items/" + strconv.Itoa(int(last_pid)) + ".txt"
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			if strings.Contains(err.Error(), "cannot assign requested address") {
				log.Println("Retrying, Error occured: [ERR500] ", err)
			}
			log.Println("Retrying, Error occured: ", err)

			continue
		}

		req.Header = http.Header{}
		req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
		req.Header.Add("accept-language", "en-US,en;q=0.9")
		req.Header.Add("cache-control", "no-cache")
		req.Header.Add("cookie", fmt.Sprintf("_vinted_fr_session=%s", session))
		req.Header.Add("pragma", "no-cache")
		req.Header.Add("priority", "u=0, i")
		req.Header.Add("sec-ch-ua", "\"Google Chrome\";v=\"125\", \"Chromium\";v=\"125\", \"Not.A/Brand\";v=\"24\"")
		req.Header.Add("sec-ch-ua-mobile", "?0")
		req.Header.Add("sec-ch-ua-platform", "\"Windows\"")
		req.Header.Add("sec-fetch-dest", "document")
		req.Header.Add("sec-fetch-mode", "navigate")
		req.Header.Add("sec-fetch-site", "none")
		req.Header.Add("sec-fetch-user", "?1")
		req.Header.Add("upgrade-insecure-requests", "1")
		req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36")

		resp, err := (*client.TlsClient).Do(req)
		if err != nil {
			log.Println("Retrying, Error occured: ", err)
			retry_count += 1
			time.Sleep(1 * time.Second)
			continue
		}
		if resp.StatusCode != 200 {
			//monitor.rotate_proxy(*client.TlsClient)
			log.Printf("[%d] not found, retrying", last_pid)
			time.Sleep(1 * time.Second)
			continue
		}
		if len(resp.Header) == 0 {
			session, err = monitor.Get_session_cookie(&client)
			if err != nil {
				log.Println("Error switching session")
			}
			println("Session")
			continue
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Retrying, Error occured: ", err)
			resp.Body.Close()
			continue
		}

		var data types.Response
		err = json.Unmarshal(body, &data)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		resp.Body.Close()
		data.Item.StringTime = time.Now().Format(time.RFC3339)
		monitor.Latest_channel <- data.Item
		last_pid = int(global_pid_list.get_new_pid())

		retry_count = 0
		continue
		//ping item

	}
}

func (l *Latest_sku) get_new_pid() int64 {
	newPid := atomic.AddInt64(&l.Latest_sku, 1)
	// Print the new PID
	fmt.Println("New PID:", newPid)
	// Return the new PID
	return newPid
}

func NewClient(_url string) (*Client, error) {
	options := Options{
		settings: []tls_client.HttpClientOption{
			tls_client.WithTimeoutSeconds(15),
			tls_client.WithClientProfile(profiles.Chrome_124),
			tls_client.WithNotFollowRedirects(),
		},
	}
	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options.settings...)
	if err != nil {
		return nil, err
	}

	return &Client{TlsClient: &client, url: _url}, nil
}

func (m *Latest_Sku_Monitor) Start_monitor() {
	log.Println("starting monitor")
	latestSku := &Latest_sku{
		Latest_sku: 0, // Initial value for Latest_sku
	}
	//log.Println("Starting a new batch at sku, ", latest_sku)
	client, err := NewClient("https://www.vinted.com")
	if err != nil {
		println("Error Creating client: ", err)
	}
	//latestSku.LatestMux.Lock()
	latestSku.Latest_sku = m.Get_latest_sku(client, m.Session) + 800
	//latestSku.LatestMux.Unlock()

	for i := 0; i < 300; i++ {

		go func() {
			client, err := NewClient("https://www.vinted.com")
			if err != nil {
				println("Error Creating client: ", err)
			}
			Make_request(int(latestSku.Latest_sku)+i, *client, m, latestSku)
		}()

	}

}

func (m *Latest_Sku_Monitor) rotate_proxy(client tls_client.HttpClient) {

	randomInt := rand.Intn(len(m.Proxies))
	err := client.SetProxy(m.Proxies[randomInt])
	if err != nil {
		log.Println("error while rotating proxies:", err)
	}
	time.Sleep(1 * time.Second)

}
