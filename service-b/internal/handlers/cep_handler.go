package handlers

import (
	"lab-cloud-run/internal/services"
	"lab-cloud-run/internal/utils"
	"net/http"
)

func CEPHandler(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Path[1:]
	err := utils.ValidateCEP(cep)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	location, err, status := services.GetLocationFromCEP(cep)
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}

	tempJSON, err := services.GetTemperature(location.Localidade)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(tempJSON))
}
