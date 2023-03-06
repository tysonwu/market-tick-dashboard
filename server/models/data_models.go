package models

import "time"

type Tick struct {
	Symbol   string    `json:"Symbol"`
	Price    float64   `json:"Price"`
	Exchange string    `json:"Exchange"`
	Time     time.Time `json:"Time"`
}

type BidAskTick struct {
	Symbol   string    `json:"Symbol"`
	BidPrice float64   `json:"Bid"`
	AskPrice float64   `json:"Ask"`
	Exchange string    `json:"Exchange"`
	Time     time.Time `json:"Time"`
}

type Client struct {
	TickChan       chan *Tick
	BidAskTickChan chan *BidAskTick
}
