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
	go exchanges.StartTickStreams(client)
	go exchanges.StartBidAskStreams(client)

	// Continuously update the Redis key-value with the latest tick data received
	go func() {
		for tick := range client.TickChan {
			// log.Println(*tick)
			dbClient.UpdateLatestTick(tick)
		}
	}()

	// not putting it into goroutine to prevent main program stops
	for tick := range client.BidAskTickChan {
		// log.Println(*tick)
		dbClient.UpdateLatestBidAskTick(tick)
	}
}
