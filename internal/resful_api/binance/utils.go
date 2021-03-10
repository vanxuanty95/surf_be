package binance

const (
	GetAggregateEndPoint = "{{host}}/api/v3/aggTrades?symbol={{symbol}}&limit={{limit}}"
	OrderEndPoint        = "{{host}}/api/v3/order?symbol={{symbol}}&side={{side}}}&type={{type}}{{params}}&newClientOrderId={{transactionID}}&timestamp={{timestamp}}&signature={{signature}}"
)

const (
	LimitPrams           = "&timeInForce={{timeInForce}}&quantity={{quantity}}&price={{price}}"
	MarketPrams          = "&quantity={{quantity}}"
	StopLossPrams        = "&quantity={{quantity}}&stopPrice={{stopPrice}}"
	StopLossLimitPrams   = "&timeInForce={{timeInForce}}&quantity={{quantity}}&price={{price}}&stopPrice={{stopPrice}}"
	TakeProfitPrams      = "&quantity={{quantity}}&stopPrice={{stopPrice}}"
	TakeProfitLimitPrams = "&timeInForce={{timeInForce}}&quantity={{quantity}}&price={{price}}&stopPrice={{stopPrice}}"
	LimitMakerPrams      = "&quantity={{quantity}}&stopPrice={{stopPrice}}"
)

const (
	TimeInForceGTC = "GTC"
	TimeInForceIOC = "IOC"
	TimeInForceFOK = "FOK"
)

const (
	OrderSideBuy  = "BUY"
	OrderSideSell = "SELL"

	OrderTypeLimit           = "LIMIT"
	OrderTypeMarket          = "MARKET"
	OrderTypeStopLoss        = "STOP_LOSS"
	OrderTypeLossLimit       = "STOP_LOSS_LIMIT"
	OrderTypeTakeProfit      = "TAKE_PROFIT"
	OrderTypeTakeProfitLimit = "TAKE_PROFIT_LIMIT"
	OrderTypeLimitMaker      = "LIMIT_MAKER"
)
