package resful_api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"surf_be/internal/app/mode"
	"surf_be/internal/app/utils"
	"surf_be/internal/configuration"
)

const (
	limitAggTrades = "1"
)

type (
	BinanceRF struct {
		Config configuration.Config
	}
)

func NewBinanceRF(config configuration.Config) BinanceRF {
	return BinanceRF{
		Config: config,
	}
}

func (ws *BinanceRF) GetAggTrades(pair string, limit string) (*mode.AggregateTrade, error) {
	if len(limit) == 0 {
		limit = limitAggTrades
	}

	endpoint := fmt.Sprintf("%v/api/v3/%ss?symbol=%s&limit=%s", ws.Config.Server.Binance.Restful.URL, utils.AggTradeStreamType, pair, limit)
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
