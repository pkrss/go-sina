package main

import (
	"log"
)

func onRecvChannel(channel string, data string) {
	log.Printf("onRecvChannel [%v] [%v]\n", channel, data)
}

func main() {
	testPlugin()
	// testDirect()

	forever := make(chan bool)
	<-forever
}
