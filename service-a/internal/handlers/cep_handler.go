package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
)

type CEPRequest struct {
	CEP string `json:"cep"`
}

func CEPHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}
	defer r.Body.Close()

	var req CEPRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	match, _ := regexp.MatchString(`^\d{8}$`, req.CEP)
	if !match {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(`invalid zipcode`))
		return
	}
}
