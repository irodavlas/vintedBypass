package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"sync"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
	"github.com/vintedMonitor/types"
	"golang.org/x/exp/rand"
)

type Sku struct {
	Sku     int
	Timeout bool
	Time    string
}
type Latest_Sku_Monitor struct {
	Proxies          []string
	Latest_channel   chan Sku
	New_batch_signal chan bool
	Latest_sku       int
	Latest_batch_sku int
	LatestMux        sync.Mutex
	Session          string
}

type Client struct {
	TlsClient *tls_client.HttpClient
	url       string
}

type Options struct {
	settings []tls_client.HttpClientOption
}

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
func (m *Latest_Sku_Monitor) Get_latest_sku(client *Client, session string) int {
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

		return int(data.Items[0].ID)
	}

}

// makes the item requests in a for loop and dies when item is found
func Make_request(sku int, session string, client *Client, latest_channel chan Sku, m *Latest_Sku_Monitor, wg *sync.WaitGroup) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var try = -1
	var general_tries = 0
	for {
		general_tries++

		if m.Latest_sku > sku {
			try++

		} else if general_tries > 10 {
			log.Println("no longer monitoring sku:", sku)
			return
		} 
		if try > 1 {
			wg.Done()
		}
		
	
		randomChar1 := charset[rand.Intn(len(charset))]
		randomChar2 := charset[rand.Intn(len(charset))]
		randomChar3 := charset[rand.Intn(len(charset))]
		url := "https://www.vinted.com/api/v2/items/" + strconv.Itoa(sku) + "." + string(randomChar1) + string(randomChar2) + string(randomChar3) 
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			m.rotate_proxy(*client.TlsClient)
			general_tries += 1
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
			println("timeout")
			log.Println("Retrying, Error occured: ", err)
			continue
		}
		if resp.StatusCode != 200 {
			log.Printf("[%d] not found, retrying", sku)
			continue
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Retrying, Error occured: ", err)
			resp.Body.Close()
			continue
		}
		defer resp.Body.Close()

		var bodyMap map[string]interface{}

		// Unmarshal JSON into the map
		err = json.Unmarshal(body, &bodyMap)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		createdAtTs, ok := bodyMap["item"].(map[string]interface{})["created_at_ts"].(string)
		if !ok {
			println("created_at_ts field not found or not of type string")
			continue
		}
		log.Println(createdAtTs)
		latest_channel <- Sku{Sku: sku, Timeout: false}

		wg.Done()
		//ping item
		return
	}
}

func NewClient(_url string) (*Client, error) {
	options := Options{
		settings: []tls_client.HttpClientOption{
			tls_client.WithTimeoutSeconds(10),
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

func (m *Latest_Sku_Monitor) Start_new_batch(latest_sku int, batch_size int) {
	//log.Println("Starting a new batch at sku, ", latest_sku)

	client, err := NewClient("https://www.vinted.com")
	if err != nil {
		println("Error starting a new batch: ", err)
	}
	sess, err := m.Get_session_cookie(client)
	if err != nil {
		println("error getting the session back")
	}
	var wg sync.WaitGroup
	for i := 0; i < batch_size; i++ {

		wg.Add(1)
		go func() {

			Make_request(latest_sku+i, sess, client, m.Latest_channel, m, &wg)
		}()

	}
	wg.Wait()
	println("go routines finished ")

	println("latest sku found:", m.Latest_sku)
	m.New_batch_signal <- true
}

func (m *Latest_Sku_Monitor) rotate_proxy(client tls_client.HttpClient) {

	randomInt := rand.Intn(len(m.Proxies))
	err := client.SetProxy(m.Proxies[randomInt])
	if err != nil {
		println("error rotating proxies")
	}

}

type Item struct {
	CreatedAtTs string `json:"created_at_ts"`
}
