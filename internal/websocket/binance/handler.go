package binance

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"surf_be/internal/app/bot"
	"surf_be/internal/app/utils"
	"surf_be/internal/configuration"
	"surf_be/internal/database/mongo"
	"surf_be/internal/mode"
)

type HandlerImpl struct {
	Config    configuration.Config
	BinanceWS BinanceWS
	Bots      map[string]*bot.Bot
	CloseCh   chan bool
	Repo      Repository
}

func NewHandler(cfg configuration.Config, mongoDB *mongo.Mongo) HandlerImpl {
	repo := NewRepository(cfg, *mongoDB.Client)
	ws := BinanceWS{Config: cfg}
	ws.Init()
	return HandlerImpl{
		BinanceWS: ws,
		Bots:      make(map[string]*bot.Bot),
		CloseCh:   make(chan bool),
		Repo:      repo,
	}
}

func (hl *HandlerImpl) PushBot(b *bot.Bot) {
	b.IDSubscribe = hl.BinanceWS.Subscribe([]string{fmt.Sprintf("%v@%v", strings.ToLower(b.Pair), b.Type)})
	hl.Bots[b.ID] = b
}

func (hl *HandlerImpl) RemoveBot(ctx context.Context, id string, idSubscribed int) {
	if err := hl.Repo.DeleteBotByID(ctx, id); err != nil {
		return
	}
	hl.BinanceWS.Unsubscribe(idSubscribed)
	delete(hl.Bots, id)
}

func (hl *HandlerImpl) GetBot(id string) *bot.Bot {
	return hl.Bots[id]
}

func (hl *HandlerImpl) DistributionMessage(cancelCtx context.Context) {
	hl.BinanceWS.Start(cancelCtx)
loop:
	for {
		select {
		case message, ok := <-hl.BinanceWS.ClientResponse:
			if !ok {
				break loop
			}
			hl.handleMessage(message)
		case <-cancelCtx.Done():
			if len(hl.BinanceWS.ClientResponse) > 0 {
				hl.handleMessage(<-hl.BinanceWS.ClientResponse)
			}
			break loop
		}
	}
	log.Println("stopped binance ws handler")
	hl.RemoveAllBot(context.Background())
	log.Println("remove all bots")
	hl.CloseCh <- true
}

func (hl *HandlerImpl) RemoveAllBot(ctx context.Context) {
	for _, botRunning := range hl.Bots {
		hl.RemoveBot(ctx, botRunning.ID, botRunning.IDSubscribe)
	}
}

func (hl *HandlerImpl) handleMessage(message []byte) {
	subscribeResponse := mode.SubscribeResponse{}
	if err := json.Unmarshal(message, &subscribeResponse); err != nil {
		log.Printf("Error reading: %v", err)
		return
	}
	if subscribeResponse.ID == 0 {
		if len(hl.Bots) == 0 {
			log.Printf("don't have bot handle this message: %v", message)
			return
		}
		detailedStreamCommon := mode.DetailedStreamCommon{}
		if err := json.Unmarshal(message, &detailedStreamCommon); err != nil {
			log.Printf("Error reading: %v", err)
			return
		}
		switch detailedStreamCommon.EventType {
		case utils.AggTradeStreamType:
			aggregateTradeModel := mode.AggregateTrade{}
			if err := json.Unmarshal(message, &aggregateTradeModel); err != nil {
				log.Printf("Error reading: %v", err)
				return
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

func (hl *HandlerImpl) isBelongTo(bot bot.Bot, detailedStreamCommon mode.DetailedStreamCommon) bool {
	return bot.Pair == detailedStreamCommon.Symbol && bot.Type == detailedStreamCommon.EventType
}
