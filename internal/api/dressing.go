package handler

import (
	Error "back/internal/err"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type NationalityResponse struct {
	Name    string `json:"name"`
	Country []struct {
		CountryID   string  `json:"country_id"`
		Probability float64 `json:"probability"`
	} `json:"country"`
}

type AgeResponse struct {
	Name  string `json:"name"`
	Age   *int   `json:"age"`
	Count int    `json:"count"`
}
type GenderResponse struct {
	Count       int     `json:"count"`
  Name        string  `json:"name"`
  Gender      string  `json:"gender"`
  Probability float64 `json:"probability"`
}



func GetNationality(name string) string{
	get := fmt.Sprintf("https://api.nationalize.io/?name=%s", name)

	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(get)
	Error.GetErr(err)
	
	defer resp.Body.Close()

	var dateRes NationalityResponse

	err = json.NewDecoder(resp.Body).Decode(&dateRes)
	Error.GetErr(err)

	return country(dateRes)
}

func ApiGet(address string) (*http.Response, error) {
	date, err := http.Get(address)
	Error.GetErr(err)

	return date, err
}

func country(nation NationalityResponse) string{
	var max float64
	var res string
	
	for _, val := range nation.Country{
		if val.Probability > max {
			max = val.Probability
			res = val.CountryID
		}
	}

	return res
}

func GetAge(name string) int {
	url := fmt.Sprintf("https://api.agify.io/?name=%s", name)

	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	Error.GetErr(err)
	
	defer resp.Body.Close()

	var age AgeResponse
	
	err = json.NewDecoder(resp.Body).Decode(&age)
	Error.GetErr(err)

	if age.Age == nil {
		return 0
	}

	return *age.Age
}


func GetGender(name string) string{
	get := fmt.Sprintf("https://api.genderize.io/?name=%s", name)

	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(get)
	Error.GetErr(err)
	
	defer resp.Body.Close()

	var gender GenderResponse

	err = json.NewDecoder(resp.Body).Decode(&gender)
  Error.GetErr(err)

	return gender.Gender
}