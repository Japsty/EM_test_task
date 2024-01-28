package server

import (
	"log"
	"net/http"
	"time"
)

type Client struct {
	Client *http.Client
}

func NewClient() *Client {
	return &Client{
		Client: &http.Client{
			Timeout: time.Minute * 3,
		},
	}
}

// SendRequest - функция делающая get запрос к стороннему сервису и возвращающая ответ от сервиса
func (c *Client) SendRequest(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return nil, err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return nil, err
	}

	return resp, nil
}
