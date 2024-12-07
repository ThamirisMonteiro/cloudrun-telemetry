package services

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type ViaCEPResponse struct {
	Localidade string `json:"localidade"`
	Estado     string `json:"estado"`
}

func GetLocationFromCEP(cep string) (*ViaCEPResponse, error, int) {
	url := "https://viacep.com.br/ws/" + cep + "/json/"

	resp, err := http.Get(url)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, err, http.StatusInternalServerError
	}

	var response ViaCEPResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	if response.Localidade == "" {
		return nil, errors.New("invalid zipcode"), http.StatusNotFound
	}

	return &response, nil, http.StatusOK
}
