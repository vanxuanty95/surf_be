package bot

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"surf_be/internal/mode"
	"time"
)

const (
	timeCalculateAveragePriceInSecond = 15.0
)

func (bot *Bot) HandleMessage(message mode.AggregateTrade) {
	now := time.Unix(0, message.EventTime*int64(time.Millisecond))

	currentPrice, err := strconv.ParseFloat(message.Price, 32)
	if err != nil {
		log.Printf("error parse float: %v", err)
		return
	}

	bot.sellOrBuy(currentPrice)

	differencePriceLatest := math.Abs(currentPrice - bot.PreviousPrice)
	differenceTimeLatest := now.Sub(bot.LatestGetTime).Seconds()
	bot.AveragePricePerSecond = bot.calculateAveragePricePerSecond(differencePriceLatest, differenceTimeLatest)
	bot.PercentSell = math.Round(bot.AveragePricePerSecond/currentPrice*1000) / 1000 * 15

	fmt.Print("\033[G\033[K")
	fmt.Printf("current  price: %v\n", currentPrice)
	fmt.Print("\033[G\033[K")
	fmt.Printf("previous price: %v\n", bot.PreviousPrice)
	fmt.Print("\033[G\033[K")
	fmt.Printf("difference price: %v\n", differencePriceLatest)
	fmt.Print("\033[G\033[K")
	fmt.Printf("difference time: %v\n", differenceTimeLatest)
	fmt.Print("\033[G\033[K")
	fmt.Printf("price per second: %v\n", differencePriceLatest/differenceTimeLatest)
	fmt.Print("\033[G\033[K")
	fmt.Printf("average price per second: %v\n", bot.AveragePricePerSecond)
	fmt.Print("\033[A")
	fmt.Print("\033[A")
	fmt.Print("\033[A")
	fmt.Print("\033[A")
	fmt.Print("\033[A")
	fmt.Print("\033[A")

	bot.CurrentData = fmt.Sprintf("current_price: %v", currentPrice)

	bot.PreviousPrice = currentPrice
	bot.LatestGetTime = now
}

func (bot *Bot) calculateAveragePricePerSecond(differencePriceLatest float64, differenceTimeLatest float64) float64 {
	if !bot.LatestGetTime.IsZero() || bot.PreviousPrice != 0 {

		totalDifferenceTimeInSlice := 0.0
		totalDifferencePriceInSlice := 0.0
		for _, priceTime := range bot.AveragePricePerSecondSlice {
			totalDifferenceTimeInSlice = totalDifferenceTimeInSlice + priceTime.differenceTime
			totalDifferencePriceInSlice = totalDifferencePriceInSlice + priceTime.differencePrice
		}

		if differenceTimeLatest+totalDifferenceTimeInSlice > timeCalculateAveragePriceInSecond {
			nextIndex := 0
			totalDifferenceTimeInSliceToRemove := 0.0
			for idx, priceTime := range bot.AveragePricePerSecondSlice {
				totalDifferenceTimeInSliceToRemove = totalDifferenceTimeInSliceToRemove + priceTime.differenceTime
				if totalDifferenceTimeInSliceToRemove >= differenceTimeLatest {
					nextIndex = idx + 1
					break
				}
			}

			bot.AveragePricePerSecondSlice = bot.AveragePricePerSecondSlice[nextIndex:]

			return (totalDifferencePriceInSlice + differencePriceLatest) / (totalDifferenceTimeInSlice + differenceTimeLatest)
		} else {
			if differencePriceLatest != 0 || differenceTimeLatest != 0 {
				bot.AveragePricePerSecondSlice = append(bot.AveragePricePerSecondSlice,
					PriceTime{
						differencePrice: differencePriceLatest,
						differenceTime:  differenceTimeLatest},
				)
			}
		}

	}
	return bot.AveragePricePerSecond
}

func (bot *Bot) sellOrBuy(currentPrice float64) {
	if bot.CurrentPrice != 0 && bot.PercentSell != 0 {
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
	} else if bot.PercentSell != 0 {
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

func (bot *Bot) GetCurrentData() string {
	return bot.CurrentData
}
