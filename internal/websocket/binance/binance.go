package binance

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"math/rand"
	"strconv"
	"surf_be/internal/configuration"
	"time"
)

const (
	durationDefaultWaitingMessageInSecond = 15
)

type (
	BinanceWS struct {
		Config         configuration.Config
		Connection     *websocket.Conn
		ClientResponse chan []byte
		SubscribedMap  map[int]Subscribed
	}
	Subscribed struct {
		ID          int
		Subscribe   string
		Unsubscribe string
	}
)

func (ws *BinanceWS) Init() {
	ws.initWSConnection()
	ws.initSubscribedMap()
	go ws.pushMessageToChannel()
}

func (ws *BinanceWS) initWSConnection() {
	c, _, err := websocket.DefaultDialer.Dial(ws.Config.Server.Binance.WebSocket.WSURL, nil)
	if err != nil {
		panic(err)
	}
	ws.Connection = c
}

func (ws *BinanceWS) initSubscribedMap() {
	ws.SubscribedMap = make(map[int]Subscribed, ws.Config.Server.Binance.WebSocket.LimitRequest)
	ws.ClientResponse = make(chan []byte)
}

func (ws *BinanceWS) pushMessageToChannel() {
	defer func() {
		err := ws.Connection.Close()
		if err != nil {
			log.Printf("error close websocket connection: %v", err)
		}
		close(ws.ClientResponse)
	}()
	tickerGetMessageTimeout := time.NewTicker(durationDefaultWaitingMessageInSecond * time.Second)
	go ws.countingTicker(tickerGetMessageTimeout)
loop:
	for {
		_, message, err := ws.Connection.ReadMessage()
		tickerGetMessageTimeout.Reset(durationDefaultWaitingMessageInSecond * time.Second)
		if err != nil || err == io.EOF {
			log.Printf("Error reading: %v", err)
			log.Printf("re-init connection")
			ws.initWSConnection()
			ws.reSubscribe()
			goto loop
		}
		ws.ClientResponse <- message
	}
}

func (ws *BinanceWS) countingTicker(ticker *time.Ticker) {
	for {
		select {
		case t := <-ticker.C:
			if len(ws.SubscribedMap) > 0 {
				log.Printf("cannot receive new message at: %v", t)
				log.Printf("re-init connection")
				ws.initWSConnection()
				ws.reSubscribe()
			}
		}
	}
}

func (ws *BinanceWS) Subscribe(pairsStr []string) int {
	id := ws.generateSubscribeString(pairsStr)
	ws.sentAMessage(ws.SubscribedMap[id].Subscribe)
	return id
}

func (ws *BinanceWS) reSubscribe() {
	for _, subscribed := range ws.SubscribedMap {
		ws.sentAMessage(subscribed.Subscribe)
	}
}

func (ws *BinanceWS) sentAMessage(message string) {
	if err := ws.Connection.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		log.Printf("write/subcribe message err: %v", err)
		ws.Connection.Close()
	}
}

func (ws *BinanceWS) generateSubscribeString(pairsStr []string) int {
	id := ws.generateID()
	pairJson, _ := json.Marshal(pairsStr)
	pairStr := fmt.Sprint(string(pairJson))

	ws.SubscribedMap[id] = Subscribed{
		ID:          id,
		Subscribe:   fmt.Sprintf("{\"method\": \"SUBSCRIBE\",\"params\": %v,\"id\": %s}", pairStr, strconv.Itoa(id)),
		Unsubscribe: fmt.Sprintf("{\"method\": \"UNSUBSCRIBE\",\"params\": %v,\"id\": %s}", pairStr, strconv.Itoa(id)),
	}
	return id
}

func (ws *BinanceWS) generateID() int {
	id := 0
	rand.Seed(time.Now().UnixNano())
	for {
		id = rand.Intn(100)
		if _, ok := ws.SubscribedMap[id]; !ok {
			return id
		}
	}
}

func (ws *BinanceWS) Unsubscribe(id int) {
	ws.sentAMessage(ws.SubscribedMap[id].Unsubscribe)
	delete(ws.SubscribedMap, id)
}
