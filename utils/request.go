package utils

import (
	"sync"
  "strings"
  "log"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
)

type Latest_Sku_Monitor struct {
  Latest_channel chan string 
  Latest_sku string
  LatestMux sync.Mutex
}

type Client struct {
	TlsClient *tls_client.HttpClient
	url       string
}

type Options struct {
	settings []tls_client.HttpClientOption
}


func get_session_cookie(session_client *Client) (*string, error){
  
		req, err := http.NewRequest(http.MethodHead, "", nil)
		if err != nil {
	    return nil, err	
			
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
	    return nil, err	
    }

		resp.Body.Close()

		cookie := resp.Header["Set-Cookie"]
		session := extractSessionCookie(cookie)
		if *session == "" {
	    return nil, err	
    }

		log.Printf("---- [%d] Fetched Session cookie----")

		return session, nil

	

}
func extractSessionCookie(cookies []string) *string {
	for _, cookie := range cookies {
		if strings.Contains(cookie, "_vinted_fr_session") {
      return &strings.Split(strings.Split(cookie, "_vinted_fr_session=")[1], ";")[0]
		}
	}
	return nil
}
func get_latest_skue(Sku_channel chan string){
}
func monitor_start(sku_to_monitor string){}

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

