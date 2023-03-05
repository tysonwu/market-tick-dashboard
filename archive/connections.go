package exchanges

import (
	"server/ticks"
)

type TickDataSource interface {
	Connect() error
	Subscribe(symbol string) error
	Receive() (*ticks.Tick, error)
}

type TickDataSources struct {
	binance TickDataSource
	// coinbase TickDataSource
}

func (t *TickDataSources) Connect() error {
	// connect to all data sources
	if err := t.binance.Connect(); err != nil {
		return err
	}
	// if err := t.coinbase.Connect(); err != nil {
	// 	return err
	// }
	// ...
	return nil
}

func (t *TickDataSources) Subscribe(symbol string) error {
	// subscribe to tick data for a symbol on all data sources
	if err := t.binance.Subscribe(symbol); err != nil {
		return err
	}
	// if err := t.coinbase.Subscribe(symbol); err != nil {
	// 	return err
	// }
	// ...
	return nil
}

func (t *TickDataSources) Receive() (*ticks.Tick, error) {
	// receive a tick update from the first data source that has new data
	for {
		tick, err := t.binance.Receive()
		if err == nil {
			return tick, nil
		}
		// tick, err = t.coinbase.Receive()
		// if err == nil {
		// 	return tick, nil
		// }
		// ...
	}
}
