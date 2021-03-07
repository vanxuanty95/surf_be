package main

import (
	"fmt"
	"log"
	"net/http"
	"surf_be/internal/configuration"
	binanceService "surf_be/internal/resful_api/binance"
	"surf_be/internal/resful_api/pfit_mgmt"
	binanceWS "surf_be/internal/websocket/binance"
)

func main() {
	env, cfgPath, err := configuration.ParseFlags()
	if err != nil {
		log.Fatal(err)
	}
	cfg, err := configuration.NewConfig(env, cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	binanceWSHL, binanceSV := initBotService(*cfg)
	initRestfulService(*cfg, binanceWSHL, binanceSV)

	stop := make(chan bool)
	<-stop
}

func initBotService(config configuration.Config) (*binanceWS.HandlerImpl, *binanceService.Service) {
	wsHandler := binanceWS.NewHandler(config)
	go wsHandler.DistributionMessage()

	binanceRF := binanceService.NewService(config)
	return &wsHandler, &binanceRF
}

func initRestfulService(config configuration.Config, binanceWSHL *binanceWS.HandlerImpl, binanceSV *binanceService.Service) {
	router := pfit_mgmt.InitPfitMgmtRouter(config, binanceWSHL, binanceSV)
	err := http.ListenAndServe(fmt.Sprintf(":%v", config.Server.PfitMgmt.APIPort), &router)
	if err != nil {
		log.Fatal(err)
	}
}
