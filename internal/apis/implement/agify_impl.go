package implement

import (
	"EM_test_task/pkg/server"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
	url := fmt.Sprintf(agify + name)

	client := server.NewClient()

	resp, err := client.SendRequest(url)
	if err != nil {
		log.Printf("Error making Agify request: %v", err)
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

//func (s *AgifyService) GetAge(name string, ch chan int) {
//	url := fmt.Sprintf("https://api.agify.io/?name=%v", name)
//
//	client := server.NewClient()
//
//	resp, err := client.SendRequest(url)
//	if err != nil {
//		log.Printf("Error making Agify request: %v", err)
//	}
//	defer resp.Body.Close()
//
//	var agifyResponse AgifyResponse
//
//	if resp.StatusCode == http.StatusOK {
//		body, err := ioutil.ReadAll(resp.Body)
//		if err != nil {
//			log.Printf("Error reading the response body:%v", err)
//		}
//		json.Unmarshal(body, &agifyResponse)
//
//		ch <- agifyResponse.Age
//	} else {
//		fmt.Println("Error: ", resp.Status)
//	}
//}
