package exchanges

import (
	"encoding/json"
	"fmt"
	"server/models"
	"strconv"
	"strings"
	"time"

	"github.com/Kucoin/kucoin-go-sdk"
	"github.com/spf13/viper"
)

type KucoinStreamConfig struct {
	Symbols []string `mapstructure:"kucoin"`
}

func makeKucoinSymbolMap(symbols []string) models.SymbolMap {
	m := models.SymbolMap{}
	for _, s := range symbols {
		m[s] = StandardizeKucoinSymbol(s)
	}
	return m
}

func StartKucoinPriceStreams(client *models.Client) {
	var streamConfig KucoinStreamConfig
	err := viper.UnmarshalKey("subscriptions", &streamConfig)
	if err != nil {
		fmt.Println("error in reading config")
		return
	}

	kucoinSymbolMap := makeKucoinSymbolMap(streamConfig.Symbols)

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

	// chname := "/market/ticker:all"
	// /market/ticker:BTC-USDT,ETH-USDT,...
	chname := fmt.Sprintf("/market/ticker:%s", strings.Join(streamConfig.Symbols[:], ","))
	ch := kucoin.NewSubscribeMessage(chname, false)
	if err := c.Subscribe(ch); err != nil {
		fmt.Println(err)
	}

	var JsonMsg struct {
		BestAsk string `json:"bestAsk"`
		BestBid string `json:"bestBid"`
		Price   string `json:"price"`
		Time    int64  `json:"time"`
	}

	for {
		select {
		case msg := <-tickChan:
			err := json.Unmarshal(msg.RawData, &JsonMsg)
			if err != nil {
				fmt.Println(err)
				// return
			}

			symbol := strings.TrimPrefix(msg.Topic, "/market/ticker:")
			price, err := strconv.ParseFloat(JsonMsg.Price, 64)
			if err != nil {
				fmt.Println(err)
				// return
			}
			bestAsk, err := strconv.ParseFloat(JsonMsg.BestAsk, 64)
			if err != nil {
				fmt.Println(err)
				// return
			}
			bestBid, err := strconv.ParseFloat(JsonMsg.BestBid, 64)
			if err != nil {
				fmt.Println(err)
				// return
			}

			// fmt.Println(symbol, price, bestAsk, bestBid, time.Unix(0, JsonMsg.Time*int64(time.Millisecond)))
			// remember: if either one of the channel did not get its message consumed somewhere else,
			// it will block the execution!

			tick := &models.Tick{
				ExchangeSymbol: symbol,
				StandardSymbol: kucoinSymbolMap[symbol],
				Price:          price,
				Exchange:       "kucoin",
				Time:           time.Unix(0, JsonMsg.Time*int64(time.Millisecond)),
			}
			client.TickChan <- tick

			bidAskTick := &models.BidAskTick{
				ExchangeSymbol: symbol,
				StandardSymbol: kucoinSymbolMap[symbol],
				BidPrice:       bestBid,
				AskPrice:       bestAsk,
				Exchange:       "kucoin",
				Time:           time.Unix(0, JsonMsg.Time*int64(time.Millisecond)),
			}
			client.BidAskTickChan <- bidAskTick

		case err := <-errorChan:
			fmt.Println(err)
		}
	}
}
