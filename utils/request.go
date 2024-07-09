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
)

const BATCH_SIZE = 500

type Latest_Sku_Monitor struct {
	Latest_channel   chan string
	New_batch_signal chan bool
	Latest_sku       string
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
func (m *Latest_Sku_Monitor) Get_latest_sku(client *Client, session string, Sku_channel chan string) {
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
		json.Unmarshal(body, &data)

		Sku_channel <- strconv.FormatInt(data.Items[0].ID, 10)
		return
	}

}

// makes the item requests in a for loop and dies when item is found
func Make_request(latest_sku string, sku string, url string, session string, client *Client, m *Latest_Sku_Monitor) {
	for {

		url := url + sku + ".txt"
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
			continue
		}
		if resp.StatusCode != 200 {
			log.Println("[%s] not found, retrying", sku)
			continue
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Retrying, Error occured: ", err)
			resp.Body.Close()
			continue
		}
		defer resp.Body.Close()

		n_to_Reach_for_new_batch, err := strconv.Atoi(m.Latest_sku)
		if err != nil {
			log.Println("Error translating the sku to int")
		}
		n_to_Reach_for_new_batch += BATCH_SIZE / 2

		if sku > n_to_Reach_for_new_batch {

		}

	}
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

func (m *Latest_Sku_Monitor) Start_new_batch() {

}
