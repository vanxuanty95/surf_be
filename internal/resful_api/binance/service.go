package binance

import (
	"encoding/json"
	"errors"
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

func (ws *Service) Oder(pair string, limit string) (*mode.AggregateTrade, error) {
	if len(limit) == 0 {
		limit = limitAggTrades
	}

	endpoint := buildLimitOrder()
	response, err := http.Get(buildLimitOrder)
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

func (ws *Service) buildLimitOrder(symbol, side, orderType, timeInForce, quantity, price, transactionID, signature string) string {

	params := strings.NewReplacer("{{timeInForce}}", timeInForce,
		"{{quantity}}", quantity,
		"{{price}}", price).Replace(LimitPrams)

	endpoint := strings.NewReplacer("{{host}}", ws.Config.Server.Binance.Restful.URL,
		"{{symbol}}", symbol,
		"{{side}}", side,
		"{{type}}", orderType,
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
