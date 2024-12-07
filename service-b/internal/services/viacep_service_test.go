package services

import (
	"testing"
)

func TestGetLocationFromCEP(t *testing.T) {
	tests := []struct {
		name    string
		cep     string
		result  *ViaCEPResponse
		wantErr bool
		status  int
	}{
		{
			name: "Valid CEP",
			cep:  "80035050",
			result: &ViaCEPResponse{
				Localidade: "Curitiba",
				Estado:     "Paraná",
			},
			wantErr: false,
			status:  200,
		},
		{
			name:    "Invalid CEP",
			cep:     "1234567a",
			result:  nil,
			wantErr: false,
			status:  500,
		},
		{
			name:    "Not found CEP",
			cep:     "00000000",
			result:  nil,
			wantErr: true,
			status:  404,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			location, err, status := GetLocationFromCEP(tt.cep)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLocationFromCEP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if location != nil {
				if location.Localidade != tt.result.Localidade || location.Estado != tt.result.Estado {
					t.Errorf("GetLocationFromCEP() location = %v, want %v, estado = %v, want %v",
						location.Localidade, tt.result.Localidade, location.Estado, tt.result.Estado)
				}
			}

			if status != tt.status {
				t.Errorf("GetLocationFromCEP() status = %v, want %v", status, tt.status)
			}
		})
	}
}
