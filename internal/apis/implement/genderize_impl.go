package implement

import (
	"EM_test_task/pkg/server"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
	url := fmt.Sprintf("https://api.genderize.io/?name=%v", name)

	client := server.NewClient()

	resp, err := client.SendRequest(url)
	if err != nil {
		log.Printf("Error making Genderize request: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	var genderizeResponce GenderizeResponce

	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading the response body:%v", err)
			return "", err
		}
		json.Unmarshal(body, &genderizeResponce)

		return genderizeResponce.Gender, nil
	} else {
		fmt.Println("Error: ", resp.Status)
	}
	return "", err
}
