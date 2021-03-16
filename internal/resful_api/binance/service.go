package binance

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"surf_be/internal/configuration"
	"surf_be/internal/mode"
	"time"
)

const (
	limitAggTrades = "1"
)

type (
	Service struct {
		Config configuration.Config
	}
)

func NewService(config configuration.Config) Service {
	return Service{
		Config: config,
	}
}

func (ws *Service) GetAggTrades(pair string, limit string) (*mode.AggregateTrade, error) {
	if len(limit) == 0 {
		limit = limitAggTrades
	}

	endpoint := strings.NewReplacer("{{host}}", ws.Config.Server.Binance.Restful.URL,
		"{{symbol}}", pair,
		"{{limit}}", limit).Replace(GetAggregateEndPoint)
	response, err := http.Get(endpoint)
	if err != nil {
		log.Printf("The HTTP request failed with error %s\n", err)
		return nil, err
	}
	defer func() {
		if err = response.Body.Close(); err != nil {
			log.Printf("close API error %s\n", err)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(response.Status)
	}

	var aggregateTrades []mode.AggregateTrade
	if err = json.NewDecoder(response.Body).Decode(&aggregateTrades); err != nil {
		return nil, err
	}

	if aggregateTrades == nil {
		return nil, errors.New("response empty")
	}

	aggregateTrades[0].Symbol = pair

	return &aggregateTrades[0], nil
}

func (ws *Service) Oder(pair, side, orderType, transactionID, signature string, quantity, price float64) (*mode.OrderResponse, error) {
	var endpoint string

	switch orderType {
	case OrderTypeLimit:
		endpoint = ws.buildLimitOrder(pair, side, transactionID, signature, quantity, price)
	default:
		return nil, errors.New("order type is not exist ")
	}

	response, err := http.Get(endpoint)
	if err != nil {
		log.Printf("The HTTP request failed with error %s\n", err)
		return nil, err
	}
	defer func() {
		if err = response.Body.Close(); err != nil {
			log.Printf("close API error %s\n", err)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(response.Status)
	}

	var errorRes mode.ErrorResponse
	if err = json.NewDecoder(response.Body).Decode(&errorRes); err != nil {
		return nil, err
	}

	if errorRes.Code != 0 {
		return nil, errors.New(errorRes.Message)
	}

	var res mode.OrderResponse
	if err = json.NewDecoder(response.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (ws *Service) buildLimitOrder(symbol, side, transactionID, signature string, quantity, price float64) string {

	params := strings.NewReplacer("{{timeInForce}}", TimeInForceGTC,
		"{{quantity}}", fmt.Sprintf("%f", quantity),
		"{{price}}", fmt.Sprintf("%f", price)).Replace(LimitPrams)

	endpoint := strings.NewReplacer("{{host}}", ws.Config.Server.Binance.Restful.URL,
		"{{symbol}}", symbol,
		"{{side}}", side,
		"{{type}}", OrderTypeLimit,
		"{{params}}", params,
		"{{transactionID}}", transactionID,
		"{{timestamp}}", ws.generateTimestamp(),
		"{{signature}}", signature,
	).Replace(OrderEndPoint)
	return endpoint
}

func (ws *Service) generateTimestamp() string {
	return strconv.FormatInt(time.Now().UnixNano()/1000000, 10)
}
