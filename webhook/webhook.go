package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/vintedMonitor/types"
)

type (
	Footer struct {
		Text    string `json:"text"`
		IconURL string `json:"icon_url"`
	}

	Field struct {
		Name   string `json:"name"`
		Value  string `json:"value"`
		Inline bool   `json:"inline"`
	}

	Thumbnail struct {
		URL string `json:"url"`
	}
	Image struct {
		URL string `json:"url"`
	}

	Embed struct {
		Title       string    `json:"title"`
		Color       int       `json:"color"`
		URL         string    `json:"url"`
		Image       Image     `json:"image"`
		Thumbnail   Thumbnail `json:"thumbnail"`
		Description string    `json:"description"`
		Footer      Footer    `json:"footer"`
		Fields      []Field   `json:"fields"`
	}

	Webhook struct {
		Content string  `json:"content"`
		Embeds  []Embed `json:"embeds"`
	}
)

func (e *Embed) SetImage(u string) {
	e.Image = Image{URL: u}
}
func (e *Embed) SetTitle(title string) {
	e.Title = title
}

func (e *Embed) SetColor(color int) {
	e.Color = color
}

func (e *Embed) SetThumbnail(u string) {
	e.Thumbnail = Thumbnail{URL: u}
}

func (e *Embed) SetDescription(description string) {
	e.Description = description
}

func (e *Embed) SetFooter(text, icon string) {
	e.Footer = Footer{Text: text, IconURL: icon}
}

func (e *Embed) SetURL(url string) {
	e.URL = url
}
func (e *Embed) AddField(name, value string, inline bool) {
	e.Fields = append(e.Fields, Field{Name: name, Value: value, Inline: inline})
}

func (w *Webhook) SetContent(content string) {
	w.Content = content
}

func (w *Webhook) AddEmbed(e Embed) {
	w.Embeds = append(w.Embeds, e)
}

func (w *Webhook) Encode() ([]byte, error) {
	return json.Marshal(w)
}

func (w *Webhook) Send(u string) error {
	payload, err := w.Encode()
	if err != nil {
		return err
	}

	for {
		req, err := http.NewRequest("POST", u, bytes.NewReader(payload))

		if err != nil {
			return err
		}

		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != 429 {
			return nil
		}

		time.Sleep(time.Second * 5)
	}
}

type Message struct {
	Item     types.ItemDetails
	Currency string
	Webhook  string
}

func Send_webhook(u string, currency string, prod types.ItemDetails) {
	webhook := &Webhook{}
	parsedTime, err := time.Parse(time.RFC3339, prod.StringTime)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return
	}
	unixTimestamp := parsedTime.Unix()
	// Create an embed
	embed := Embed{}
	Title := prod.Title
	Description := fmt.Sprintf("Posted: <t:%d:R>", unixTimestamp)
	Price := currency + prod.Price.Amount
	Photo := prod.Photos[0].URL
	Store := prod.Country
	Size := prod.SizeTitle
	Condition := prod.Status
	BrandTitle := prod.Title

	webhook.SetContent("")
	embed.SetTitle(Title)
	embed.SetURL(prod.URL)
	embed.SetDescription(Description)
	embed.SetColor(5549236)
	embed.SetImage(Photo)
	embed.SetFooter(fmt.Sprintf("NRC BOT • %s", prod.StringTime), "https://cdn.discordapp.com/attachments/1220385098541826161/1249380232163496047/NRCUK.png?ex=66671783&is=6665c603&hm=7c3a76b0ca29a947a14ca25f2791dfa8e447bbc3dd3fddfd9f065379d821bf69&")

	embed.AddField("💷 Item Price", Price, false)

	embed.AddField("Store", Store, false)
	embed.AddField("📏 Size", Size, true)
	embed.AddField("🏅 Condition", Condition, true)
	embed.AddField("👕 Brand", BrandTitle, true)

	// Add the embed to the webhook
	webhook.AddEmbed(embed)

	// Send the webhook to the specified URL
	webhookURL := u
	err = webhook.Send(webhookURL)
	if err != nil {
		fmt.Println("Error sending webhook:", err)
	} else {
		fmt.Println("Webhook sent successfully")
	}
}
