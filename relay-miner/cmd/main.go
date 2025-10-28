package main

import (
	"fmt"

	"github.com/SEG-UNIBE/artio-insight/relay-miner/pkg/miner"
)

/*
main example function for just fetching the data from the relays
*/
func main() {
	relay := "relay.artiostr.ch"

	miner.GetNip11(fmt.Sprintf("https://%v/", relay))
	miner.GetRelayList(fmt.Sprintf("wss://%v/", relay))

}
