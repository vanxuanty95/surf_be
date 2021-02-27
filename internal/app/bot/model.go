package bot

import (
	"time"
)

type Bot struct {
	ID            int
	StartTime     time.Time
	Duration      time.Duration
	Pair          string
	Access        string
	BuyInPrice    float64
	BuyInQuantity float64
	CurrentPrice  float64
	SellPrice     float64
	SellTime      time.Time
	Quantity      float64
	StopChannel   chan bool
	Type          string
	IDSubscribe   int
	PercentBuy    float64
	PercentSell   float64
	Budget        float64
}
