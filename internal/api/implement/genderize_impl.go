package implement

import (
	"EM_test_task/pkg/client"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type GenderizeService struct{}

type GenderizeResponce struct {
	Count       int     `json:"count"`
	Name        string  `json:"name"`
	Gender      string  `json:"gender"`
	Probability float64 `json:"float64"`
}

// GetGender - метод для похода в API Genderuze, получает пол по имени
func (s *GenderizeService) GetGender(name string) (string, error) {
	gender := os.Getenv("GENDERIZE_URL")
	urlQuery := os.Getenv("URL_QUERY")
	url := fmt.Sprintf(gender + urlQuery + name)

	newClient := client.NewClient()

	resp, err := newClient.GetAPIResponseByURL(url)
	if err != nil {
		log.Printf("Error making Genderize request: %v", err)
		return "", err
	}

	body, err := newClient.ReadResponseBody(resp)
	if err != nil {
		log.Printf("Error reading Agify response body: %v", err)
		return "", err
	}

	var genderizeResponse GenderizeResponce
	err = json.Unmarshal(body, &genderizeResponse)
	if err != nil {
		log.Printf("Error unmarshaling json: %v", err)
		return "", err
	}
	return genderizeResponse.Gender, nil
}
