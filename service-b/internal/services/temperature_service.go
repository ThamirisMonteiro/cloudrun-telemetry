package services

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
)

type WeatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

func GetTemperature(location string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar o arquivo .env")
	}

	// Pega a chave de API do ambiente
	weatherAPIKey := os.Getenv("WEATHER_API_KEY")
	if weatherAPIKey == "" {
		log.Fatal("A chave de API não foi definida em WEATHER_API_KEY")
	}

	url := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s", weatherAPIKey, location)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var response WeatherAPIResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	tempF := response.Current.TempC*1.8 + 32
	tempK := response.Current.TempC + 273.15

	tempC := math.Round(response.Current.TempC*10) / 10
	tempF = math.Round(tempF*10) / 10
	tempK = math.Round(tempK*10) / 10

	temperature := map[string]float64{
		"temp_C": tempC,
		"temp_F": tempF,
		"temp_K": tempK,
	}

	temperatureJSON, err := json.Marshal(temperature)
	if err != nil {
		return "", err
	}

	return string(temperatureJSON), nil
}
