package main

import (
	"github.com/pkrss/go-sina/plugins/quote"
)

func testPlugin() {
	var p = onRecvChannel

	f := quote.Start(p)

	f("sh601988,RB0")
}
