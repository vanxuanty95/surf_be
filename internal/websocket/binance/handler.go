package binance

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"surf_be/internal/app/bot"
	"surf_be/internal/app/utils"
	"surf_be/internal/configuration"
	"surf_be/internal/mode"
)

type HandlerImpl struct {
	Config    configuration.Config
	BinanceWS BinanceWS
	Bots      map[string]*bot.Bot
}

func NewHandler(cfg configuration.Config) HandlerImpl {
	ws := BinanceWS{Config: cfg}
	ws.Init()
	return HandlerImpl{
		BinanceWS: ws,
		Bots:      make(map[string]*bot.Bot),
	}
}

func (hl *HandlerImpl) PushBot(b *bot.Bot) {
	b.IDSubscribe = hl.BinanceWS.Subscribe([]string{fmt.Sprintf("%v@%v", strings.ToLower(b.Pair), b.Type)})
	hl.Bots[b.ID] = b
}

func (hl *HandlerImpl) RemoveBot(id string) {
	delete(hl.Bots, id)
}

func (hl *HandlerImpl) DistributionMessage() {
	for {
		select {
		case message := <-hl.BinanceWS.ClientResponse:
			subscribeResponse := mode.SubscribeResponse{}
			if err := json.Unmarshal(message, &subscribeResponse); err != nil {
				log.Printf("Error reading: %v", err)
				continue
			}
			if subscribeResponse.ID == 0 {
				if len(hl.Bots) == 0 {
					log.Printf("don't have bot handle this message: %v", message)
					continue
				}
				detailedStreamCommon := mode.DetailedStreamCommon{}
				if err := json.Unmarshal(message, &detailedStreamCommon); err != nil {
					log.Printf("Error reading: %v", err)
					continue
				}
				switch detailedStreamCommon.EventType {
				case utils.AggTradeStreamType:
					aggregateTradeModel := mode.AggregateTrade{}
					if err := json.Unmarshal(message, &aggregateTradeModel); err != nil {
						log.Printf("Error reading: %v", err)
						continue
					}
					for _, b := range hl.Bots {
						if hl.isBelongTo(*b, detailedStreamCommon) {
							b.HandleMessage(aggregateTradeModel)
						}
					}
				default:
					log.Printf("out of type handler: %s", detailedStreamCommon.EventType)
					log.Printf("message: %s", message)
				}
			} else {
				log.Printf("subscribe successed!: %v", subscribeResponse.ID)
			}
		}
	}
}

func (hl *HandlerImpl) isBelongTo(bot bot.Bot, detailedStreamCommon mode.DetailedStreamCommon) bool {
	return bot.Pair == detailedStreamCommon.Symbol && bot.Type == detailedStreamCommon.EventType
}
