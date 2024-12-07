package services

import (
	"strings"
	"testing"
)

func TestGetTemperature(t *testing.T) {
	temperature, err := GetTemperature("Curitiba")
	if err != nil {
		t.Errorf("GetTemperature() error = %v, wantErr %v", err, false)
	}
	if !strings.Contains(temperature, "temp_C") || !strings.Contains(temperature, "temp_F") || !strings.Contains(temperature, "temp_K") {
		t.Errorf("GetTemperature() temperature = %v, want %v", temperature, "temp_C, temp_F, temp_K")
	}
}
