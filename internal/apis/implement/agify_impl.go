package implement

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type AgifyService struct{}

type AgifyResponse struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
}

func (s *AgifyService) GetAge(name string) (int, error) {
	url := fmt.Sprintf("https://api.agify.io/?name=%v", name)

	client := http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating request:%v", err)
		return 0, err
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request:%v", err)
		return 0, err
	}
	defer resp.Body.Close()

	var agifyResponse AgifyResponse

	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading the response body:%v", err)
			return 0, err
		}
		json.Unmarshal(body, &agifyResponse)

		return agifyResponse.Age, nil
	} else {
		fmt.Println("Error: ", resp.Status)
	}
	return 0, err
}
