package implement

import (
	"EM_test_task/pkg/server"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type NationalizeService struct{}

type NationalizeResponse struct {
	Count     int           `json:"count"`
	Name      string        `json:"name"`
	Countries []CountryInfo `json:"country"`
}

type CountryInfo struct {
	CountryName string  `json:"country_id"`
	Probability float64 `json:"probability"`
}

// GetNationality - метод для похода в API Nationalize, получает национальность по фамилии
// (в тз было указано имя в запросе, но в доке апишки написана фамилия) )
func (s *NationalizeService) GetNationality(surname string) (string, error) {
	url := fmt.Sprintf("https://api.nationalize.io/?name=%v", surname)

	client := server.NewClient()

	resp, err := client.SendRequest(url)
	if err != nil {
		log.Printf("Error making Nationalize request: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	var nationalizeResponce NationalizeResponse

	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading the response body:%v", err)
			return "", err
		}
		json.Unmarshal(body, &nationalizeResponce)

		var maxProb float64
		var foundConuntry string
		for _, country := range nationalizeResponce.Countries {
			if country.Probability > maxProb {
				maxProb = country.Probability
				foundConuntry = country.CountryName
			}
		}

		return foundConuntry, nil
	} else {
		fmt.Println("Error: ", resp.Status)
	}
	return "", err
}
