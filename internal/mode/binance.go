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

	ErrorResponse struct {
		Code    int    `json:"code"`
		Message string `json:"msg"`
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

	OrderResponse struct {
		Symbol              string `json:"symbol"`
		OrderId             int64  `json:"orderId"`
		OrderListId         string `json:"orderListId"`
		ClientOrderId       int    `json:"clientOrderId"`
		TransactTime        string `json:"transactTime"`
		Price               string `json:"price"`
		OrigQty             int    `json:"origQty"`
		ExecutedQty         int    `json:"executedQty"`
		CummulativeQuoteQty int64  `json:"cummulativeQuoteQty"`
		Status              bool   `json:"status"`
		TimeInForce         bool   `json:"timeInForce"`
		Type                bool   `json:"type"`
		Side                bool   `json:"side"`
		Fills               bool   `json:"fills"`
	}
)
