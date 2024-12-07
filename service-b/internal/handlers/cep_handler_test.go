package handlers

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCEPHandler_ValidCEP_RealAPI(t *testing.T) {
	cep := "80035050"

	req, err := http.NewRequest("GET", "/"+cep, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(CEPHandler)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Esperado status 200 OK, mas obteve %v", rr.Code)
	}

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]float64
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, response, "temp_C")
	assert.Contains(t, response, "temp_F")
	assert.Contains(t, response, "temp_K")

	assert.Equal(t, float64(int(response["temp_C"]*10))/10, response["temp_C"])
	assert.Equal(t, float64(int(response["temp_F"]*10))/10, response["temp_F"])
	assert.Equal(t, float64(int(response["temp_K"]*10))/10, response["temp_K"])
}

func TestCEPHandler_InvalidCEP_RealAPI(t *testing.T) {
	cep := "invalidCEP"
	req, err := http.NewRequest("GET", "/"+cep, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CEPHandler)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnprocessableEntity {
		t.Errorf("Esperado status 422, mas obteve %v", rr.Code)
	}

	expected := "invalid zipcode"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("Esperado corpo da resposta %v, mas obteve %v", expected, rr.Body.String())
	}
}
