package exchanges

import (
	"fmt"
	"server/models"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/spf13/viper"
)

type StreamConfig struct {
	Symbols []string `mapstructure:"binance"`
}

func StartTickStreams(client *models.Client) {
	var streamConfig StreamConfig
	err := viper.UnmarshalKey("subscriptions", &streamConfig)
	if err != nil {
		fmt.Println("error in reading config")
		return
	}

	wsAggTradeHandler := func(event *binance.WsAggTradeEvent) {
		price, err := strconv.ParseFloat(event.Price, 64)
		if err != nil {
			fmt.Println(err)
			return
		}
		tick := &models.Tick{
			Symbol:   event.Symbol,
			Price:    price,
			Exchange: "binance",
			Time:     time.Unix(0, event.Time*int64(time.Millisecond)),
		}
		client.TickChan <- tick
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

func StartBidAskStreams(client *models.Client) {
	var streamConfig StreamConfig
	err := viper.UnmarshalKey("subscriptions", &streamConfig)
	if err != nil {
		fmt.Println("error in reading config")
		return
	}

	wsBookTickerHandler := func(event *binance.WsBookTickerEvent) {
		bestBidPrice, err := strconv.ParseFloat(event.BestBidPrice, 64)
		if err != nil {
			fmt.Println(err)
			return
		}
		bestAskPrice, err := strconv.ParseFloat(event.BestAskPrice, 64)
		if err != nil {
			fmt.Println(err)
			return
		}
		bidAskTick := &models.BidAskTick{
			Symbol:   event.Symbol,
			BidPrice: bestBidPrice,
			AskPrice: bestAskPrice,
			Exchange: "binance",
			Time:     time.Now(), // binance API did not give the time field in this ws
		}
		client.BidAskTickChan <- bidAskTick
	}
	errHandler := func(err error) {
		fmt.Println(err)
	}

	doneC, _, err := binance.WsCombinedBookTickerServe(streamConfig.Symbols, wsBookTickerHandler, errHandler)
	if err != nil {
		fmt.Println(err)
		return
	}
	<-doneC
}
