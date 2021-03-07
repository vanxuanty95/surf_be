package bot

import (
	"time"
)

type Bot struct {
	ID          int
	StartTime   time.Time
	Duration    time.Duration
	Pair        string
	Access      string
	Type        string
	IDSubscribe int
	StopChannel chan bool

	BuyInPrice    float64
	BuyInQuantity float64
	CurrentPrice  float64
	PreviousPrice float64
	LatestGetTime time.Time
	SellPrice     float64
	SellTime      time.Time
	Quantity      float64

	PercentBuy  float64
	PercentSell float64
	Budget      float64

	AveragePricePerSecondSlice []PriceTime
	AveragePricePerSecond      float64

	CurrentData string
}

type PriceTime struct {
	differencePrice float64
	differenceTime  float64
}
