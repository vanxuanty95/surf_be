package pfit_mgmt

import (
	"github.com/gorilla/mux"
	"surf_be/internal/configuration"
	"surf_be/internal/database/redis"
	"surf_be/internal/resful_api/binance"
	binanceWS "surf_be/internal/websocket/binance"
)

const (
	PostMethod = "POST"
	GetMethod  = "GET"

	PathToPing      = "/ping"
	PathToLogout    = "/logout"
	PathToLogin     = "/login"
	GetTransactions = "/transactions"
	GetBotStatus    = "/bot/status/{pair}"
)

func InitPfitMgmtRouter(config configuration.Config, binanceHL *binanceWS.HandlerImpl, binanceSV *binance.Service) mux.Router {
	redisDB := redis.Redis{Config: config}
	redisDB.Init()

	InitValidationService()

	service := NewService(config, redisDB, binanceSV, binanceHL)
	handler := NewHandler(config, service)

	r := mux.NewRouter()
	r.Use(Authorize(redisDB))

	r.Path(PathToLogin).Methods(PostMethod).HandlerFunc(handler.Login)
	r.Path(GetTransactions).Methods(PostMethod).HandlerFunc(handler.GetTransactions)
	r.Path(GetBotStatus).Methods(GetMethod).HandlerFunc(handler.GetBotStatus)

	return *r
}
