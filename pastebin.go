package pastebin

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

const (
	scrapeURL = "https://scrape.pastebin.com/api_scraping.php?limit=250"
	pasteURL  = "https://scrape.pastebin.com/api_scrape_item.php?i="
)

// Paste is a struct that represents a paste object from Pastebin's API.  I
// excluded full_url, scrape_url and size, to reduce space usage, as they can be
// derived from the paste key.
type Paste struct {
	Date      string `json:"date"`
	Key       string `json:"key"`
	Expire    string `json:"expire"`
	Title     string `json:"title"`
	Syntax    string `json:"syntax"`
	User      string `json:"user"`
	FullURL   string `json:"full_url"`
	ScrapeURL string `json:"scrape_url"`
	Size      string `json:"size"`
}

type Client struct {
	client *http.Client
}

func NewClient() *Client {
	httpClient := newHTTPClient()
	return &Client{
		httpClient,
	}
}

func newHTTPClient() *http.Client {
	var transport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}

	var client = &http.Client{
		Timeout:   time.Second * 10,
		Transport: transport,
	}

	return client
}

func (c *Client) LatestPastes() ([]Paste, error) {
	resp, err := c.client.Get(scrapeURL)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	pastes := make([]Paste, 0)

	err = json.NewDecoder(resp.Body).Decode(&pastes)
	if err != nil {
		return nil, err
	}

	return pastes, nil
}

func (c *Client) GetPaste(key string) (string, error) {
	resp, err := c.client.Get(pasteURL + key)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	return string(content), nil
}
