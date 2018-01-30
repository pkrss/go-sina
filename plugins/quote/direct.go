package quote

import (
	"github.com/pkrss/go-sina/quote"
)

func startDirectMode(p func(channel string, data string)) func(string) {
	quote.OnPublishChannel = p

	go quote.Start()

	return quote.NotifyStkChanged
}
