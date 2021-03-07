package pfit_mgmt

import (
	"encoding/json"
	"net/http"
)

func DecodeBody(r *http.Request, requestMode interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(&requestMode); err != nil {
		return err
	}
	return nil
}

//======================================
type LoginRequest struct {
	Email     string `json:"email"`
	APIKey    string `json:"api_key"`
	SecretKey string `json:"secret_key"`
}
