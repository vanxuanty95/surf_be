package pfit_mgmt

import (
	"github.com/gorilla/mux"
	"surf_be/internal/configuration"
	"surf_be/internal/database/redis"
)

const (
	PostMethod = "POST"
)

func InitPfitMgmtRouter(config configuration.Config) mux.Router {
	redisDB := redis.Redis{Config: config}
	redisDB.Init()

	handler := NewHandler(config, redisDB)

	r := mux.NewRouter()
	r.Path("/login").Methods(PostMethod).HandlerFunc(handler.Login)

	return *r
}
