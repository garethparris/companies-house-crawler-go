package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
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

	if rest_api_key == "" {
		fmt.Print("No API key found")
		os.Exit(1)
	}

	company_number := "10833732"

	officers, err := getOfficersForCompany(company_number)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	for i := 0; i < len(officers); i++ {

		fmt.Println()

		dumpOfficer(officers[i])

		var officer_id = extractOfficerId(officers[i])
		if officer_id == "" {
			fmt.Println("No officer ID found")
			continue
		}

		appointments, err := getAppointmentsForOfficer(officer_id)

		if err != nil {
			fmt.Print(err.Error())
			continue
		}

		for j := 0; j < len(appointments); j++ {
			dumpAppointment(appointments[j])

			officers2, err := getOfficersForCompany(appointments[j].AppointedTo.Number)

			if err != nil {
				fmt.Print(err.Error())
				continue
			}

			for k := 0; k < len(officers2); k++ {
				dumpOfficer(officers[k])
			}
		}
	}
}

func getOfficersForCompany(company_number string) ([]OfficerSummary, error) {
	var bearer = "Basic " + base64.StdEncoding.EncodeToString([]byte(rest_api_key))
	url := host + "/company/" + company_number + "/officers"

	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Add("Authorization", bearer)

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	responseData, err := io.ReadAll(io.Reader(response.Body))
	if err != nil {
		return nil, err
	}

	var responseObject OfficerResponse
	json.Unmarshal(responseData, &responseObject)

	return responseObject.Officers, nil
}

func getAppointmentsForOfficer(officer_id string) ([]AppointmentSummary, error) {
	var bearer = "Basic " + base64.StdEncoding.EncodeToString([]byte(rest_api_key))
	url := host + "/officers/" + officer_id + "/appointments"

	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Add("Authorization", bearer)

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	responseData, err := io.ReadAll(io.Reader(response.Body))
	if err != nil {
		return nil, err
	}

	var responseObject AppointmentResponse
	json.Unmarshal(responseData, &responseObject)

	return responseObject.Appointments, nil
}

func extractOfficerId(officerSummary OfficerSummary) string {
	parts := strings.Split(officerSummary.Links.Officer.Appointments, "/")
	if len(parts) >= 3 {
		return parts[2]
	}
	return ""
}

func dumpOfficer(officer OfficerSummary) {

	fmt.Println(officer.Name)

	// fmt.Println(responseObject.Officers[i].Nationality)
	// fmt.Println(responseObject.Officers[i].OfficerRole)
	// fmt.Println(
	// 	strconv.Itoa(responseObject.Officers[i].DateOfBirth.Month) + "-" +
	// 		strconv.Itoa(responseObject.Officers[i].DateOfBirth.Year))
	// fmt.Println(responseObject.Officers[i].Links.Self)
}

func dumpAppointment(appointment AppointmentSummary) {
	fmt.Println(appointment.AppointedTo.Name)
	/*
		fmt.Println(appointment.AppointedOn)
		fmt.Println(appointment.AppointedTo.Number)
		fmt.Println(appointment.AppointedTo.Status)
	*/
}
