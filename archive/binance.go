package exchanges

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"server/ticks"

	"github.com/gorilla/websocket"
)

type BinanceDataSource struct {
	conn      *websocket.Conn
	tickChan  chan *ticks.Tick
	closeChan chan struct{}
}

func (b *BinanceDataSource) Connect() error {
	dialer := websocket.Dialer{
		Proxy:             http.ProxyFromEnvironment,
		HandshakeTimeout:  45 * time.Second,
		EnableCompression: false,
	}
	conn, _, err := dialer.Dial("wss://stream.binance.com:9443/ws", nil)
	if err != nil {
		return fmt.Errorf("failed to connect to Binance: %v", err)
	}

	b.conn = conn
	b.tickChan = make(chan *ticks.Tick)
	b.closeChan = make(chan struct{})

	return nil
}

func (b *BinanceDataSource) Subscribe(symbol string) error {
	subscribeMsg := map[string]interface{}{
		"method": "SUBSCRIBE",
		"params": []string{fmt.Sprintf("%s@aggTrade", symbol)},
		"id":     1,
	}

	err := b.conn.WriteJSON(subscribeMsg)
	if err != nil {
		return fmt.Errorf("failed to subscribe to %s on Binance: %v", symbol, err)
	}

	return nil
}

func (b *BinanceDataSource) Receive() (*ticks.Tick, error) {
	var tickMsg struct {
		Symbol string `json:"s"`
		Price  string `json:"p"`
		Time   string `json:"E"`
	}
	// var tickMsg map[string]interface{}
	var jsonMsg map[string]interface{}

	_, msg, err := b.conn.ReadMessage()
	if err != nil {
		return nil, fmt.Errorf("failed to read message from Binance: %v", err)
	}

	err = json.Unmarshal(msg, &jsonMsg)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	if _, ok := jsonMsg["result"]; ok {
		// first message is a map with key result indicating the start of stream
		// call the method again recursively
		return b.Receive()
	}

	err = json.Unmarshal(msg, &tickMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal tick msg: %v", err)
	}

	// convert price and time type
	// price, err := strconv.ParseFloat(tickMsg.Price, 64)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to parse price: %v", err)
	// }

	// timeMillis, err := strconv.ParseInt(tickMsg.Time, 10, 64)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to parse time: %v", err)
	// }
	// timeUnix := timeMillis * int64(time.Millisecond)

	fmt.Println(tickMsg)
	tick := &ticks.Tick{
		// Symbol:   tickMsg.Symbol,
		// Price:    tickMsg.Price,
		// Exchange: "Binance",
		// Time:     time.Unix(0, timeUnix),
	}
	return tick, nil
}
