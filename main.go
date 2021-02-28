package main

import (
	"fmt"
	"log"
	"strconv"
	"surf_be/internal/app/bot"
	"surf_be/internal/app/utils"
	"surf_be/internal/configuration"
	"surf_be/internal/resful_api"
	"surf_be/internal/websocket"
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
	wsHandler := websocket.NewHandler(*cfg)
	go wsHandler.DistributionMessage()

	binanceRF := resful_api.NewBinanceRF(*cfg)

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
		PercentSell:   0.01,
		Budget:        0,
	}

	wsHandler.PushBot(&BTCBot)

	stop := make(chan bool)
	<-stop
}
