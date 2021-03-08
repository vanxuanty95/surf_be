package pfit_mgmt

import (
	"encoding/json"
	"errors"
	"net/http"
)

func DecodeBody(r *http.Request, requestMode interface{}) error {
	if r.Body == nil {
		return errors.New("request body is empty")
	}
	if err := json.NewDecoder(r.Body).Decode(&requestMode); err != nil {
		return err
	}
	return nil
}

//======================================
type LoginRequest struct {
	Email     string `json:"email" validate:"required,email"`
	APIKey    string `json:"api_key" validate:"required"`
	SecretKey string `json:"secret_key" validate:"required"`
}

type StartBotRequest struct {
	Access string `json:"access" validate:"required"`
	Pair   string `json:"pair" validate:"required"`
}
