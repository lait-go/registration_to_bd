package handler

import (
	Error "back/internal/err"
	LogWork "back/internal/utils/log"
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
	Age   *int   `json:"age"`   // 👈 указатель — это правильно
	Count int    `json:"count"`
}

type GenderResponse struct {
	Count       int     `json:"count"`
  Name        string  `json:"name"`
  Gender      string  `json:"gender"`
  Probability float64 `json:"probability"`
}

var sharedClient = &http.Client{Timeout: 3 * time.Second}


func SafeGetAge(name string) (result int) {
	defer func() {
		if r := recover(); r != nil {
			LogWork.LogInit().Debug(fmt.Sprintf("panic в SafeGetAge: %v", r))
			result = 0
		}
	}()

	logger := LogWork.LogInit()
	url := fmt.Sprintf("https://api.agify.io/?name=%s", name)
	resp, err := sharedClient.Get(url)
	if Error.GetErr(err) {
		logger.Debug(fmt.Sprintf("Ошибка запроса возраста: %v", err))
		return 0
	}
	defer resp.Body.Close()

	var ageResp AgeResponse
	if err := json.NewDecoder(resp.Body).Decode(&ageResp); Error.GetErr(err) {
		logger.Debug(fmt.Sprintf("Ошибка разбора возраста: %v", err))
		return 0
	}
	if ageResp.Age == nil {
		logger.Debug("Поле Age отсутствует в ответе API")
		return 0
	}
	return *ageResp.Age
}

func SafeGetGender(name string) string {
	logger := LogWork.LogInit()
	url := fmt.Sprintf("http://api.genderize.io/?name=%s", name)
	resp, err := sharedClient.Get(url)
	if Error.GetErr(err) {
		logger.Debug(fmt.Sprintf("Ошибка запроса пола: %v", err))
		return ""
	}
	defer resp.Body.Close()

	var genderResp GenderResponse
	if err := json.NewDecoder(resp.Body).Decode(&genderResp); Error.GetErr(err) {
		logger.Debug(fmt.Sprintf("Ошибка разбора пола: %v", err))
		return ""
	}
	return genderResp.Gender
}

func SafeGetNationality(name string) string {
	logger := LogWork.LogInit()
	url := fmt.Sprintf("http://api.nationalize.io/?name=%s", name)
	resp, err := sharedClient.Get(url)
	if Error.GetErr(err) {
		logger.Debug(fmt.Sprintf("Ошибка запроса национальности: %v", err))
		return ""
	}
	defer resp.Body.Close()

	var natResp NationalityResponse
	if err := json.NewDecoder(resp.Body).Decode(&natResp); Error.GetErr(err) {
		logger.Debug(fmt.Sprintf("Ошибка разбора национальности: %v", err))
		return ""
	}
	return extractBestCountry(natResp)
}

func extractBestCountry(n NationalityResponse) string {
	var max float64
	var result string
	for _, c := range n.Country {
		if c.Probability > max {
			max = c.Probability
			result = c.CountryID
		}
	}
	return result
}