package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"surf_be/internal/configuration"
	"surf_be/internal/database/mongo"
	"surf_be/internal/database/redis"
	binanceService "surf_be/internal/resful_api/binance"
	"surf_be/internal/resful_api/pfit_mgmt"
	binanceWS "surf_be/internal/websocket/binance"
	"syscall"
)

func main() {
	srvCtx, srvCancel := context.WithCancel(context.Background())

	cfg := initConfiguration()

	redisDB, mongoDB := initDatabase(*cfg)

	binanceWSHL, binanceSV := initBotService(srvCtx, *cfg, mongoDB)
	restSrv := initRestfulService(*cfg, redisDB, mongoDB, binanceWSHL, binanceSV)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	<-signals

	log.Println("shutting down http server")
	if err := restSrv.Shutdown(srvCtx); err != nil {
		log.Fatalf("http server shutdown with error: %v", err)
	}
	srvCancel()
	<-binanceWSHL.CloseCh
}

func initConfiguration() *configuration.Config {
	env, cfgPath, err := configuration.ParseFlags()
	if err != nil {
		log.Fatal(err)
	}
	cfg, err := configuration.NewConfig(env, cfgPath)
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}

func initBotService(cancelCtx context.Context, config configuration.Config, mongoDB *mongo.Mongo) (*binanceWS.HandlerImpl, *binanceService.Service) {
	wsHandler := binanceWS.NewHandler(config, mongoDB)
	go wsHandler.DistributionMessage(cancelCtx)

	binanceRF := binanceService.NewService(config)
	return &wsHandler, &binanceRF
}

func initRestfulService(config configuration.Config, redisDB *redis.Redis, mongoDB *mongo.Mongo, binanceWSHL *binanceWS.HandlerImpl, binanceSV *binanceService.Service) *http.Server {
	router := pfit_mgmt.InitPfitMgmtRouter(config, redisDB, mongoDB, binanceWSHL, binanceSV)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", config.Server.PfitMgmt.APIPort),
		Handler: &router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	return srv
}

func initDatabase(config configuration.Config) (*redis.Redis, *mongo.Mongo) {
	redisDB := redis.Redis{Config: config}
	redisDB.Init()
	mongoDB := mongo.Mongo{Config: config}
	mongoDB.Init()

	return &redisDB, &mongoDB
}
