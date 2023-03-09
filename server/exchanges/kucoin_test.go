package exchanges

import (
	"fmt"
	"log"
	"testing"

	"github.com/Kucoin/kucoin-go-sdk"
)

func TestStartKucoinPriceStreams(t *testing.T) {
	// define an array of string with name streamConfig with two values 'A' and 'B'
	var streamConfig = []string{"BTC-USDT", "ETH-USDT"}

	service := kucoin.NewApiServiceFromEnv()

	rsp, err := service.WebSocketPublicToken()
	if err != nil {
		fmt.Println(err)
	}

	tk := &kucoin.WebSocketTokenModel{}
	if err := rsp.ReadData(tk); err != nil {
		fmt.Println(err)
	}

	c := service.NewWebSocketClient(tk)

	tickChan, errorChan, err := c.Connect()
	if err != nil {
		fmt.Println(err)
	}

	for _, symbol := range streamConfig {
		log.Println(symbol)
		chname := fmt.Sprintf("/market/ticker:%s", symbol)
		ch := kucoin.NewSubscribeMessage(chname, false)
		if err := c.Subscribe(ch); err != nil {
			fmt.Println(err)
		}
	}

	for {
		select {
		case tick := <-tickChan:
			fmt.Println(tick)
		case err := <-errorChan:
			fmt.Println(err)
		}
	}

	// wsAggTradeHandler := func(event *kucoin.WsAggTradeEvent) {
	// 	price, err := strconv.ParseFloat(event.Price, 64)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return
	// 	}
	// 	tick := &models.Tick{
	// 		Symbol:   event.Symbol,
	// 		Price:    price,
	// 		Exchange: "binance",
	// 		Time:     time.Unix(0, event.Time*int64(time.Millisecond)),
	// 	}
	// 	client.TickChan <- tick
	// }

	// doneC, _, err := binance.WsCombinedAggTradeServe(streamConfig.Symbols, wsAggTradeHandler, errHandler)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// <-doneC
}
