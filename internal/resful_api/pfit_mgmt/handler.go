package pfit_mgmt

import (
	"net/http"
	"surf_be/internal/configuration"
)

type (
	Handler struct {
		Config  configuration.Config
		Service Service
	}
)

func NewHandler(config configuration.Config, service Service) Handler {
	return Handler{
		Config:  config,
		Service: service,
	}
}

func (hl *Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := LoginRequest{}

	if err := DecodeBody(r, &req); err != nil {
		WriteJSON(w)(HandleError(err, 9999))
		return
	}

	res, err := hl.Service.Login(ctx, req)
	if err != nil {
		WriteJSON(w)(HandleError(err, 9999))
		return
	}

	WriteJSON(w)(HandleSuccess(res))
}

func (hl *Handler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w)(HandleSuccess("ok"))
}
