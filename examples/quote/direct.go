package main

import (
	"github.com/pkrss/go-sina/quote"
)

func testDirect() {
	var p = onRecvChannel

	quote.OnPublishChannel = p

	go quote.Start()

	f := quote.NotifyStkChanged

	f("sh601988,RB0")
}
