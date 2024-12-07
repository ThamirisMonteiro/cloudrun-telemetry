package utils

import "testing"

func TestValidateCEP(t *testing.T) {
	tests := []struct {
		name    string
		cep     string
		wantErr bool
	}{
		{
			name:    "Valid CEP",
			cep:     "80035050",
			wantErr: false,
		},
		{
			name:    "Invalid CEP",
			cep:     "1234567a",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateCEP(tt.cep); (err != nil) != tt.wantErr {
				t.Errorf("ValidateCEP() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
