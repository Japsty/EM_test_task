package implement

import (
	"EM_test_task/pkg/client"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type NationalizeService struct{}

type nationalizeResponse struct {
	Count     int           `json:"count"`
	Name      string        `json:"name"`
	Countries []countryInfo `json:"country"`
}

type countryInfo struct {
	CountryName string  `json:"country_id"`
	Probability float64 `json:"probability"`
}

// GetNationality - метод для похода в API Nationalize, получает национальность по фамилии
// (в тз было указано имя в запросе, но в доке API написана фамилия)
func (s *NationalizeService) GetNationality(surname string) (string, error) {
	nation := os.Getenv("NATIONALIZE_URL")
	urlQuery := os.Getenv("URL_QUERY")
	url := fmt.Sprintf(nation + urlQuery + surname)

	newClient := client.NewClient()

	resp, err := newClient.GetAPIResponseByURL(url)
	if err != nil {
		log.Printf("Error making Nationalize request: %v", err)
		return "", err
	}
	body, err := newClient.ReadResponseBody(resp)
	if err != nil {
		log.Printf("Error reading Agify response body: %v", err)
		return "", err
	}

	var nationalizeResponse nationalizeResponse
	err = json.Unmarshal(body, &nationalizeResponse)
	if err != nil {
		log.Printf("Error unmarshaling json: %v", err)
		return "", err
	}

	var maxProb float64
	var foundCountry string
	for _, country := range nationalizeResponse.Countries {
		if country.Probability > maxProb {
			maxProb = country.Probability
			foundCountry = country.CountryName
		}
	}
	return foundCountry, nil
}
