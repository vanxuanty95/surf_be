package pfit_mgmt

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteJSON(w http.ResponseWriter) func(resp interface{}, statusCode int) {
	return func(resp interface{}, statusCode int) {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(statusCode)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("error when encode response: %v", err)
		}
	}
}

func HandleError(err error, statusCode int) (ErrorResponse, int) {
	return ErrorResponse{
		Code:  statusCode,
		Error: err.Error(),
	}, http.StatusInternalServerError
}

func HandleSuccess(rep interface{}) (SuccessResponse, int) {
	return SuccessResponse{
		Code: 0000,
		Data: rep,
	}, http.StatusOK
}

type (
	ErrorResponse struct {
		Code  int    `json:"code"`
		Error string `json:"error"`
	}
	SuccessResponse struct {
		Code int         `json:"code"`
		Data interface{} `json:"data"`
	}
)

//======================================
type LoginResponse struct {
	EmailWithEnv string `json:"email_with_env"`
}

type StartBotResponse struct {
	BotID string `json:"bot_id"`
}

type GetBotStatusResponse struct {
	IsStarted    bool   `json:"is_started"`
	BotID        string `json:"bot_id"`
	CurrentPrice string `json:"current_price"`
}
