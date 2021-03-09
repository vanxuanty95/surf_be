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

func (hl *Handler) StartBot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := StartBotRequest{}

	if err := DecodeBody(r, &req); err != nil {
		WriteJSON(w)(HandleError(err, 9999))
		return
	}

	res, err := hl.Service.StartBot(ctx, req)
	if err != nil {
		WriteJSON(w)(HandleError(err, 9999))
		return
	}

	WriteJSON(w)(HandleSuccess(res))
}

func (hl *Handler) GetBotStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := GetBotStatusRequest{}

	if err := DecodeBody(r, &req); err != nil {
		WriteJSON(w)(HandleError(err, 9999))
		return
	}

	res, err := hl.Service.GetGetBotStatus(ctx, req)
	if err != nil {
		WriteJSON(w)(HandleError(err, 9999))
		return
	}

	WriteJSON(w)(HandleSuccess(res))
}

func (hl *Handler) StopBot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := StopBotRequest{}

	if err := DecodeBody(r, &req); err != nil {
		WriteJSON(w)(HandleError(err, 9999))
		return
	}

	err := hl.Service.StopBot(ctx, req)
	if err != nil {
		WriteJSON(w)(HandleError(err, 9999))
		return
	}

	WriteJSON(w)(HandleSuccess(nil))
}
