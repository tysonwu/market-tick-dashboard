package ticks

import "time"

type Tick struct {
	Symbol   string    `json:"Symbol"`
	Price    float64   `json:"Price"`
	Exchange string    `json:"Exchange"`
	Time     time.Time `json:"Time"`
}

type TickClient struct {
	TickChan chan *Tick
}
