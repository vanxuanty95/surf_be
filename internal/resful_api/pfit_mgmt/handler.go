package pfit_mgmt

import (
	"net/http"
	"surf_be/internal/configuration"
	"surf_be/internal/database/redis"
)

type (
	Handler struct {
		Config  configuration.Config
		RedisDB redis.Redis
	}
)

func NewHandler(config configuration.Config, redisDB redis.Redis) Handler {
	return Handler{
		Config:  config,
		RedisDB: redisDB,
	}
}

func (gl *Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := LoginRequest{}
	res := LoginResponse{}

	if err := DecodeBody(r, &req); err != nil {
		WriteJSON(w)(HandleError(err, 9999))
		return
	}

	if err := gl.RedisDB.Set(ctx, req.Email, req, 300); err != nil {
		WriteJSON(w)(HandleError(err, 9999))
		return
	}

	res.Email = req.Email

	WriteJSON(w)(HandleSuccess(res))
}
