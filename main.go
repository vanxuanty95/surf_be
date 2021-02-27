package main

import (
	"log"
	"surf_be/internal/app/bot"
	"surf_be/internal/app/utils"
	"surf_be/internal/configuration"
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
	WSHandler := websocket.NewHandler(*cfg)
	go WSHandler.DistributionMessage()

	BTCBot := bot.Bot{
		ID:            1,
		StartTime:     time.Now(),
		Duration:      2 * time.Hour,
		Pair:          "DOTUSDT",
		Access:        "DOT",
		BuyInPrice:    30.332199096679688,
		BuyInQuantity: 20,
		CurrentPrice:  30.332199096679688,
		Quantity:      20,
		StopChannel:   nil,
		Type:          utils.AggTradeStreamType,
		PercentBuy:    0.1,
		PercentSell:   1,
		Budget:        0,
	}

	WSHandler.PushBot(&BTCBot)

	stop := make(chan bool)
	<-stop
}
