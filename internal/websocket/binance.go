package websocket

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
	c, _, err := websocket.DefaultDialer.Dial(ws.Config.Server.Binance.WebSocket, nil)
	if err != nil {
		panic(err)
	}
	ws.SubscribedMap = make(map[int]Subscribed, ws.Config.Server.Binance.LimitRequest)
	ws.ClientResponse = make(chan []byte)
	ws.Connection = c
	go ws.pushMessageToChannel()
}

func (ws *BinanceWS) pushMessageToChannel() {
	defer close(ws.ClientResponse)
	for {
		_, message, err := ws.Connection.ReadMessage()
		if err != nil || err == io.EOF {
			log.Fatal("Error reading: ", err)
			break
		}
		ws.ClientResponse <- message
	}
}

func (ws *BinanceWS) Subscribe(pairsStr []string) int {
	id := ws.generateSubscribeString(pairsStr)
	if err := ws.Connection.WriteMessage(websocket.TextMessage, []byte(ws.SubscribedMap[id].Subscribe)); err != nil {
		ws.Connection.Close()
	}
	return id
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
	if err := ws.Connection.WriteMessage(websocket.TextMessage, []byte(ws.SubscribedMap[id].Unsubscribe)); err != nil {
		ws.Connection.Close()
	}
	delete(ws.SubscribedMap, id)
}
