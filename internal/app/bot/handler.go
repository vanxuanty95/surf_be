package bot

import (
	"log"
	"strconv"
	"surf_be/internal/app/mode"
	"time"
)

func (bot *Bot) Start() {
	bot.StartTime = time.Now()

loop:
	for {
		select {

		case <-time.After(bot.Duration):
			break loop
		case <-bot.StopChannel:
			break loop
		}
	}
}

func (bot *Bot) startGetMessage() {

}

func (bot *Bot) inputMessageTo() {

}

func (bot *Bot) HandleMessage(message mode.AggregateTrade) {
	currentPrice, err := strconv.ParseFloat(message.Price, 32)
	if err != nil {
		log.Printf("error parse float: %v", err)
		return
	}

	if bot.CurrentPrice != 0 {
		sellPointPrice := bot.CurrentPrice + bot.CurrentPrice/100*bot.PercentSell
		buyPointPrice := bot.CurrentPrice - bot.CurrentPrice/100*bot.PercentBuy

		if currentPrice >= sellPointPrice {
			profit := currentPrice*bot.Quantity - bot.BuyInPrice*bot.BuyInQuantity
			bot.Budget = currentPrice*bot.Quantity + bot.Budget
			bot.CurrentPrice = 0
			bot.Quantity = 0
			bot.SellPrice = currentPrice
			bot.SellTime = time.Now()

			log.Printf("in price at: %v", bot.BuyInPrice)
			log.Printf("sell at:     %v", currentPrice)
			log.Printf("profit:      %v", profit)
			log.Printf("quantity:    %v", bot.Quantity)
			log.Printf("budget:      %v", bot.Budget)
			log.Printf("================")
		}
		if bot.Budget > 0 && currentPrice <= buyPointPrice {
			quantityAdd := bot.Budget / currentPrice
			bot.Budget = 0
			bot.CurrentPrice = currentPrice
			bot.Quantity = quantityAdd

			log.Printf("in price at: %v", bot.BuyInPrice)
			log.Printf("buy at:      %v", currentPrice)
			log.Printf("quantity:    %v", bot.Quantity)
			log.Printf("budget:      %v", bot.Budget)
			log.Printf("================")
		}
	} else {
		buyPointPrice := bot.SellPrice - bot.SellPrice/100*bot.PercentBuy

		if bot.Budget > 0 && (currentPrice <= buyPointPrice ||
			(((currentPrice-bot.SellPrice)/bot.SellPrice*100) > bot.PercentBuy) && time.Now().Sub(bot.SellTime) > 15*time.Second) {
			quantityAdd := bot.Budget / currentPrice
			bot.Budget = 0
			bot.CurrentPrice = currentPrice
			bot.Quantity = quantityAdd

			log.Printf("in price at: %v", bot.BuyInPrice)
			log.Printf("buy at:      %v", currentPrice)
			log.Printf("quantity:    %v", bot.Quantity)
			log.Printf("budget:      %v", bot.Budget)
			log.Printf("================")
		}
	}
}
