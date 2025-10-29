package main

import (
	"github.com/SEG-UNIBE/artio-insight/relay-miner/pkg/miner"
)

/*
main example function for just fetching the data from the relays
*/
func main() {
	startingRelays := []string{"relay.artiostr.ch", "relay.artio.inf.unibe.ch"}
	miners := make([]*miner.RelayMiner, 0)
	for _, relay := range startingRelays {
		miners = append(miners, miner.NewMiner(relay))
	}

	for _, relay := range miners {
		relay.Load()
		relay.Stats()
	}
}
