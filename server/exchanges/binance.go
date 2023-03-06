package exchanges

import (
	"fmt"
	"server/ticks"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/spf13/viper"
)

type StreamConfig struct {
	Symbols []string `mapstructure:"binance"`
}

func Start(tickClient *ticks.TickClient) {
	var streamConfig StreamConfig
	err := viper.UnmarshalKey("subscriptions", &streamConfig)
	if err != nil {
		fmt.Println("error in reading config")
		return
	}
	fmt.Println(streamConfig.Symbols)

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

	doneC, _, err := binance.WsCombinedAggTradeServe(streamConfig.Symbols, wsAggTradeHandler, errHandler)
	if err != nil {
		fmt.Println(err)
		return
	}
	<-doneC
}
