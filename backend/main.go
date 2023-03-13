package main

import (
	"fmt"
	"log"
	"server/db"
	"server/exchanges"
	"server/models"

	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile("./config.yaml") // path from go program root
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}

func main() {
	// to prevent main thread from exiting
	done := make(chan bool)

	dbClient, err := db.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	defer dbClient.Close()

	go dbClient.PublishTick()
	go dbClient.PublishBidAskTick()

	// Start the Binance websocket stream
	client := &models.Client{
		TickChan:       make(chan *models.Tick),
		BidAskTickChan: make(chan *models.BidAskTick),
	}

	// ============BINANCE===========
	go exchanges.StartTickStreams(client)
	go exchanges.StartBidAskStreams(client)

	// Continuously update the Redis key-value with the latest tick data received
	go func() {
		for tick := range client.TickChan {
			// log.Println(*tick)
			dbClient.UpdateLatestTick(tick)
		}
	}()

	go func() {
		for tick := range client.BidAskTickChan {
			// log.Println(*tick)
			dbClient.UpdateLatestBidAskTick(tick)
		}
	}()

	// ============KUCOIN===========
	go exchanges.StartKucoinPriceStreams(client)

	// Continuously update the Redis key-value with the latest tick data received
	go func() {
		for tick := range client.TickChan {
			// log.Println(*tick)
			dbClient.UpdateLatestTick(tick)
		}
	}()

	go func() {
		for tick := range client.BidAskTickChan {
			// log.Println(*tick)
			dbClient.UpdateLatestBidAskTick(tick)
		}
	}()

	<-done
}
