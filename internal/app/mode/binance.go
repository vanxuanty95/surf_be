package mode

type (
	SubscribeResponse struct {
		Result string `json:"result"`
		ID     int    `json:"id"`
	}

	DetailedStreamCommon struct {
		EventType string `json:"e"`
		EventTime int64  `json:"E"`
		Symbol    string `json:"s"`
	}

	AggregateTrade struct {
		EventType          string `json:"e"`
		EventTime          int64  `json:"E"`
		Symbol             string `json:"s"`
		AggregateTradeID   int    `json:"a"`
		Price              string `json:"p"`
		Quantity           string `json:"q"`
		FirstTradeID       int    `json:"f"`
		LastTradeID        int    `json:"l"`
		TradeTime          int64  `json:"T"`
		IsBuyerMarketMaker bool   `json:"m"`
		Ignore             bool   `json:"M"`
	}
)
