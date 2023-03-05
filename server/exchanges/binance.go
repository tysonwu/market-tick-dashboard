package exchanges

import (
	"fmt"
	"server/ticks"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2"
)

func Start(tickClient *ticks.TickClient) {
	wsAggTradeHandler := func(event *binance.WsAggTradeEvent) {
		price, err := strconv.ParseFloat(event.Price, 64)
		if err != nil {
			fmt.Println(err)
			return
		}
		tick := &ticks.Tick{
			Symbol:   event.Symbol,
			Price:    price,
			Exchange: "binance",
			Time:     time.Unix(0, event.Time*int64(time.Millisecond)),
		}
		tickClient.TickChan <- tick
	}
	errHandler := func(err error) {
		fmt.Println(err)
	}

	symbols := []string{"btcusdt", "ethusdt", "bnbusdt"}
	doneC, _, err := binance.WsCombinedAggTradeServe(symbols, wsAggTradeHandler, errHandler)
	if err != nil {
		fmt.Println(err)
		return
	}
	<-doneC
}
