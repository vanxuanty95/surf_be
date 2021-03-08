package bot

import (
	"time"
)

type Bot struct {
	ID          string        `bson:"id"`
	StartTime   time.Time     `bson:"start_time"`
	Duration    time.Duration `bson:"duration"`
	Pair        string        `bson:"pair"`
	Access      string        `bson:"access"`
	Type        string        `bson:"type"`
	IDSubscribe int           `bson:"id_subscribe"`
	Email       string        `bson:"email"`
	StopChannel chan bool     `bson:"-"`
	CurrentData string        `bson:"-"`

	BuyInPrice    float64   `bson:"-"`
	BuyInQuantity float64   `bson:"-"`
	CurrentPrice  float64   `bson:"-"`
	PreviousPrice float64   `bson:"-"`
	LatestGetTime time.Time `bson:"-"`
	SellPrice     float64   `bson:"-"`
	SellTime      time.Time `bson:"-"`
	Quantity      float64   `bson:"-"`

	PercentBuy  float64 `bson:"-"`
	PercentSell float64 `bson:"-"`
	Budget      float64 `bson:"-"`

	AveragePricePerSecondSlice []PriceTime `bson:"-"`
	AveragePricePerSecond      float64     `bson:"-"`
}

type PriceTime struct {
	differencePrice float64
	differenceTime  float64
}
