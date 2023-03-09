package exchanges

import (
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/adshao/go-binance/v2"
)

func TestBinanceConnectivity(t *testing.T) {
	wsAggTradeHandler := func(event *binance.WsAggTradeEvent) {
		price, err := strconv.ParseFloat(event.Price, 64)
		log.Println(price)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	errHandler := func(err error) {
		fmt.Println(err)
	}

	symbols := []string{"btcusdt"}
	doneC, stopC, err := binance.WsCombinedAggTradeServe(symbols, wsAggTradeHandler, errHandler)
	if err != nil {
		fmt.Println(err)
		return
	}
	<-time.After(3 * time.Second)
	stopC <- struct{}{}
	<-doneC
	<-time.After(1 * time.Second)
}
