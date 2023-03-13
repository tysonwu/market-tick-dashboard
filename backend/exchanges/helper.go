package exchanges

import (
	"regexp"
	"strings"
)

func StandardizeBinanceSymbol(symbol string) string {
	// find base currency then add hyphens
	regex := regexp.MustCompile(`([a-zA-Z]+)*(USDT|USDC|USD|BTC$)`)
	match := regex.FindStringSubmatch(symbol)
	res := strings.ToUpper(strings.Join(match[1:], "-")) // index 0 is the original symbol
	return res
}

func StandardizeKucoinSymbol(symbol string) string {
	// find base currency then add hyphens
	// nothing to change as it is already in format
	return strings.ToUpper(symbol)
}
