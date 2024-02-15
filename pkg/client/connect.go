package client

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Client struct {
	Client http.Client
}

func NewClient() Client {
	return Client{
		Client: http.Client{
			Timeout: 1 * time.Second,
		},
	}
}

// GetAPIResponseByURL - функция делающая get запрос к стороннему сервису и возвращающая ответ от сервиса
func (c *Client) GetAPIResponseByURL(url string) (*http.Response, error) {
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

func (c *Client) ReadResponseBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading the response body:%v", err)
			return nil, err
		}
		return body, nil
	}
	log.Printf("ReadResponseBody Error, response code: %v", resp.StatusCode)
	err := fmt.Errorf("ReadResponseBody Error, response code: %v", resp.StatusCode)
	return nil, err
}
