package implement

import (
	"EM_test_task/pkg/client"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type AgifyService struct {
}

type AgifyResponse struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
}

// GetAge - метод для похода в API Agify, получает возраст по имени
func (s *AgifyService) GetAge(name string) (int, error) {
	agify := os.Getenv("AGIFY_URL")
	urlQuery := os.Getenv("URL_QUERY")
	url := fmt.Sprintf(agify + urlQuery + name)

	newClient := client.NewClient()

	resp, err := newClient.GetAPIResponseByURL(url)
	if err != nil {
		log.Printf("Error making Agify request: %v", err)
		return 0, err
	}

	body, err := newClient.ReadResponseBody(resp)
	if err != nil {
		log.Printf("Error reading Agify response body: %v", err)
		return 0, err
	}

	var agifyResponse AgifyResponse
	err = json.Unmarshal(body, &agifyResponse)
	if err != nil {
		log.Printf("Error unmarshaling json: %v", err)
		return 0, err
	}
	return agifyResponse.Age, nil
}
