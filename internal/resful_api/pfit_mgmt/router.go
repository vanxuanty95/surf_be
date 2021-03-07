package pfit_mgmt

import (
	"github.com/gorilla/mux"
	"surf_be/internal/configuration"
	"surf_be/internal/database/redis"
)

const (
	PostMethod = "POST"

	PathToPing      = "/ping"
	PathToLogout    = "/logout"
	PathToLogin     = "/login"
	GetTransactions = "/transactions"
)

func InitPfitMgmtRouter(config configuration.Config) mux.Router {
	redisDB := redis.Redis{Config: config}
	redisDB.Init()

	InitValidationService()

	service := NewService(config, redisDB)
	handler := NewHandler(config, service)

	r := mux.NewRouter()
	r.Use(Authorize(redisDB))

	r.Path(PathToLogin).Methods(PostMethod).HandlerFunc(handler.Login)
	r.Path(GetTransactions).Methods(PostMethod).HandlerFunc(handler.GetTransactions)

	return *r
}
