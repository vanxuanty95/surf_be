package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"surf_be/internal/app/bot"
	"surf_be/internal/app/utils"
	"surf_be/internal/configuration"
	binanceService "surf_be/internal/resful_api/binance"
	"surf_be/internal/resful_api/pfit_mgmt"
	binanceWS "surf_be/internal/websocket/binance"
	"time"
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

	initRestfulService(*cfg)
	//initBotService(*cfg)

	stop := make(chan bool)
	<-stop
}

func initBotService(config configuration.Config) {
	wsHandler := binanceWS.NewHandler(config)
	go wsHandler.DistributionMessage()

	binanceRF := binanceService.NewService(config)

	access := "DOT"
	excess := "USDT"
	pair := fmt.Sprintf("%v%v", access, excess)

	rspData, err := binanceRF.GetAggTrades(pair, "1")
	if err != nil {
		log.Fatal(err)
	}

	currentPrice, err := strconv.ParseFloat(rspData.Price, 32)
	if err != nil {
		log.Fatalf("error parse float: %v", err)
	}

	BTCBot := bot.Bot{
		ID:            1,
		StartTime:     time.Now(),
		Duration:      2 * time.Hour,
		Pair:          rspData.Symbol,
		Access:        access,
		BuyInPrice:    currentPrice,
		BuyInQuantity: 1,
		CurrentPrice:  currentPrice,
		Quantity:      1,
		StopChannel:   nil,
		Type:          utils.AggTradeStreamType,
		PercentBuy:    0.01,
		Budget:        0,
	}

	wsHandler.PushBot(&BTCBot)
}

func initRestfulService(config configuration.Config) {
	router := pfit_mgmt.InitPfitMgmtRouter(config)
	err := http.ListenAndServe(fmt.Sprintf(":%v", config.Server.PfitMgmt.APIPort), &router)
	if err != nil {
		log.Fatal(err)
	}
}
