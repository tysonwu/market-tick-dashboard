package models

import "time"

type Tick struct {
	StandardSymbol string    `json:"StandardSymbol"`
	ExchangeSymbol string    `json:"ExchangeSymbol"`
	Price          float64   `json:"Price"`
	Exchange       string    `json:"Exchange"`
	Time           time.Time `json:"Time"`
}

type BidAskTick struct {
	StandardSymbol string    `json:"StandardSymbol"`
	ExchangeSymbol string    `json:"ExchangeSymbol"`
	BidPrice       float64   `json:"Bid"`
	AskPrice       float64   `json:"Ask"`
	Exchange       string    `json:"Exchange"`
	Time           time.Time `json:"Time"`
}

type Client struct {
	TickChan       chan *Tick
	BidAskTickChan chan *BidAskTick
}

type SymbolMap map[string]string
