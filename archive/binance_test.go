package exchanges

import (
	"testing"
)

func TestBinanceWsConnectivity(t *testing.T) {
	binance := BinanceDataSource{}

	err := binance.Connect()
	if err != nil {
		t.Errorf("failed to connect to Binance: %v", err)
	}

	symbol := "btcusdt"
	err = binance.Subscribe(symbol)
	if err != nil {
		t.Errorf("failed to subscribe to symbol: %v", err)
	}

	for i := 0; i < 3; i++ {
		_, err := binance.Receive()
		if err != nil {
			t.Errorf("failed to receive tick data: %v", err)
		}
	}
}
