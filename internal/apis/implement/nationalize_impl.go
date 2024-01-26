package implement

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type NationalizeResponse struct {
	Count     int           `json:"count"`
	Name      string        `json:"name"`
	Countries []CountryInfo `json:"country"`
}

type CountryInfo struct {
	CountryName string  `json:"country_id"`
	Probability float64 `json:"probability"`
}

type NationalizeService struct{}

func (s *NationalizeService) GetNationality(surname string) (string, error) {
	url := fmt.Sprintf("https://api.nationalize.io/?name=%v", surname)

	client := http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating request:%v", err)
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request:%v", err)
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
