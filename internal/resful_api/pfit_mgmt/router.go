package pfit_mgmt

import (
	"github.com/gorilla/mux"
	"surf_be/internal/configuration"
	"surf_be/internal/database/mongo"
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
	StartBot        = "/bot/start"
	GetBotStatus    = "/bot/status"
)

func InitPfitMgmtRouter(config configuration.Config, redisDB *redis.Redis, mongoDB *mongo.Mongo, binanceHL *binanceWS.HandlerImpl, binanceSV *binance.Service) mux.Router {
	InitValidationService()

	repo := NewRepository(config, *mongoDB.Client)
	service := NewService(config, redisDB, repo, binanceSV, binanceHL)
	handler := NewHandler(config, service)

	r := mux.NewRouter()
	r.Use(Authorize(redisDB))

	r.Path(PathToLogin).Methods(PostMethod).HandlerFunc(handler.Login)
	r.Path(GetTransactions).Methods(PostMethod).HandlerFunc(handler.GetTransactions)
	r.Path(StartBot).Methods(PostMethod).HandlerFunc(handler.StartBot)
	r.Path(GetBotStatus).Methods(PostMethod).HandlerFunc(handler.GetBotStatus)

	return *r
}
