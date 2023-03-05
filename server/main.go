package main

import (
	"fmt"
	"log"
	"server/db"
	"server/exchanges"
	"server/ticks"

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

	// Start the Binance websocket stream
	tickClient := &ticks.TickClient{TickChan: make(chan *ticks.Tick)}
	go exchanges.Start(tickClient)

	// Continuously update the Redis key-value with the latest tick data received
	for tick := range tickClient.TickChan {
		log.Println(tick)
		dbClient.UpdateLatestTick(tick)
	}
}
