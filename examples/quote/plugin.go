package main

import (
	"github.com/pkrss/gosina/plugins/quote"
)

func testPlugin() {
	var p = onRecvChannel

	f := quote.Start(p)

	f("sh601988,RB0")
}
