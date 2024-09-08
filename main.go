package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// A Response struct to map the Entire Response
type Response struct {
	ActiveCount int              `json:"active_count"`
	ETag        string           `json:"etag"`
	Officers    []OfficerSummary `json:"items"`
}

// A struct to map our Pokemon's Species which includes it's name
type OfficerSummary struct {
	Name        string `json:"name"`
	Nationality string `json:"nationality"`
}

func main() {

	rest_api_key := os.Getenv("COMP_HOUSE_API_KEY")

	url := "https://api.company-information.service.gov.uk/company/10833732/officers"

	// Create a Bearer string by appending string access token
	var bearer = "Basic " + base64.StdEncoding.EncodeToString([]byte(rest_api_key))

	// Create a new request using http
	request, _ := http.NewRequest("GET", url, nil)

	// add authorization header to the req
	request.Header.Add("Authorization", bearer)

	// Send req using http Client
	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := io.ReadAll(io.Reader(response.Body))
	if err != nil {
		log.Fatal(err)
	}

	var responseObject Response
	json.Unmarshal(responseData, &responseObject)

	fmt.Println(responseObject.ActiveCount)
	fmt.Println(responseObject.ETag)
	fmt.Println(len(responseObject.Officers))

	for i := 0; i < len(responseObject.Officers); i++ {
		fmt.Println(responseObject.Officers[i].Name)
		fmt.Println(responseObject.Officers[i].Nationality)
	}
}
