package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	Exchange string  `json:"exchange"`
	Symbol   string  `json:"symbol"`
	Price    float64 `json:"price"`
	Time     int64   `json:"time"`
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a web socket connection
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v\n", err)
		return
	}
	defer conn.Close()

	// Loop indefinitely and send price updates
	for {
		// Simulate fetching price data from various exchanges
		// Replace this code with your own implementation to fetch real price data
		for _, exchange := range []string{"binance", "coinbase", "bitstamp"} {
			for _, symbol := range []string{"BTCUSD", "ETHUSD", "LTCUSD"} {
				price := 10000.0 + (float64(time.Now().UnixNano()%1000000)/1000000)*100.0
				message := Message{
					Exchange: exchange,
					Symbol:   symbol,
					Price:    price,
					Time:     time.Now().UnixNano(),
				}
				data, err := json.Marshal(message)
				if err != nil {
					log.Printf("Error marshaling message: %v\n", err)
					continue
				}
				if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
					log.Printf("Error sending message: %v\n", err)
					continue
				}
				log.Printf("%v\n", price)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
