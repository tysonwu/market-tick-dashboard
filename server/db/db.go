package db

import (
	"context"
	"encoding/json"
	"fmt"
	"server/ticks"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

type DbClient struct {
	client   *redis.Client
	TickChan chan *ticks.Tick
}

type RedisConfig struct {
	Address  string `mapstructure:"address"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

func NewClient() (*DbClient, error) {
	var redisConfig RedisConfig
	err := viper.UnmarshalKey("redis", &redisConfig)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(&redis.Options{
		Addr:     redisConfig.Address,
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})

	return &DbClient{
		client:   client,
		TickChan: make(chan *ticks.Tick, 1), // buffer channel with size 1
	}, nil
}

func (db *DbClient) PublishTick() {
	// loop to continuously read the chan
	for tick := range db.TickChan {
		symbol := tick.Symbol
		exchange := tick.Exchange
		jsonData, err := json.Marshal(tick)
		if err != nil {
			fmt.Println(err)
			continue
		}

		var m map[string]interface{}
		err = json.Unmarshal(jsonData, &m)
		if err != nil {
			fmt.Println(err)
			continue
		}

		keyName := fmt.Sprintf("ticks:%s:%s", symbol, exchange)
		err = db.client.HSet(context.Background(), keyName, m).Err() // set <string: hash>
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}

func (db *DbClient) UpdateLatestTick(tick *ticks.Tick) {
	// Send the latest tick data to tickChan, overwriting any existing data in the channel
	db.TickChan <- tick
}

func (db *DbClient) Close() {
	db.client.Close()
}