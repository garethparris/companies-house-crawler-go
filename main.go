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

type OfficerResponse struct {
	ActiveCount int              `json:"active_count"`
	ETag        string           `json:"etag"`
	Officers    []OfficerSummary `json:"items"`
}

type OfficerSummary struct {
	Name        string      `json:"name"`
	Nationality string      `json:"nationality"`
	Occupation  string      `json:"occupation"`
	OfficerRole string      `json:"officer_role"`
	DateOfBirth DateOfBirth `json:"date_of_birth"`
	Links       Links       `json:"links"`
}

type DateOfBirth struct {
	Day   int `json:"day"`
	Month int `json:"month"`
	Year  int `json:"year"`
}

type Links struct {
	Self    string      `json:"self"`
	Officer OfficerLink `json:"officer"`
}

type OfficerLink struct {
	Appointments string `json:"appointments"`
}

type AppointmentResponse struct {
	ActiveCount  int                  `json:"active_count"`
	ETag         string               `json:"etag"`
	Appointments []AppointmentSummary `json:"items"`
}

type AppointmentSummary struct {
	AppointedOn string      `json:"appointed_on"`
	AppointedTo AppointedTo `json:"appointed_to"`
}

type AppointedTo struct {
	Name   string `json:"company_name"`
	Number string `json:"company_number"`
	Status string `json:"company_status"`
}

const host = "https://api.company-information.service.gov.uk"

var rest_api_key string

func main() {

	rest_api_key = os.Getenv("COMP_HOUSE_API_KEY")

	company_number := "10833732"

	getOfficersForCompany(company_number)

	getAppointmentsForOfficer("KXhQ15Scka6lnwt7rRBxRcC_ueg")
}

func getOfficersForCompany(company_number string) {

	var bearer = "Basic " + base64.StdEncoding.EncodeToString([]byte(rest_api_key))
	url := host + "/company/" + company_number + "/officers"

	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Add("Authorization", bearer)

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

	var responseObject OfficerResponse
	json.Unmarshal(responseData, &responseObject)

	fmt.Println(responseObject.ActiveCount)
	fmt.Println(responseObject.ETag)
	fmt.Println(len(responseObject.Officers))

	for i := 0; i < len(responseObject.Officers); i++ {
		fmt.Println(responseObject.Officers[i].Name)
		fmt.Println(responseObject.Officers[i].Links.Officer.Appointments)
		// fmt.Println(responseObject.Officers[i].Nationality)
		// fmt.Println(responseObject.Officers[i].OfficerRole)
		// fmt.Println(
		// 	strconv.Itoa(responseObject.Officers[i].DateOfBirth.Month) + "-" +
		// 		strconv.Itoa(responseObject.Officers[i].DateOfBirth.Year))
		// fmt.Println(responseObject.Officers[i].Links.Self)
	}
}

func getAppointmentsForOfficer(officer_id string) {

	var bearer = "Basic " + base64.StdEncoding.EncodeToString([]byte(rest_api_key))
	url := host + "/officers/" + officer_id + "/appointments"

	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Add("Authorization", bearer)

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

	var responseObject AppointmentResponse
	json.Unmarshal(responseData, &responseObject)

	fmt.Println(responseObject.ActiveCount)
	fmt.Println(responseObject.ETag)
	fmt.Println(len(responseObject.Appointments))

	for i := 0; i < len(responseObject.Appointments); i++ {
		fmt.Println(responseObject.Appointments[i].AppointedOn)
		fmt.Println(responseObject.Appointments[i].AppointedTo.Name)
		fmt.Println(responseObject.Appointments[i].AppointedTo.Number)
		fmt.Println(responseObject.Appointments[i].AppointedTo.Status)
		fmt.Println()
	}
}
