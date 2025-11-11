package main

import (
	"log"

	"github.com/SEG-UNIBE/artio-insight/relay-miner/pkg/miner"
	"github.com/SEG-UNIBE/artio-insight/relay-miner/pkg/storage"
)

/*
main example function for just fetching the data from the relays
*/
func main() {
	startingRelays := []string{"wss://relay.artiostr.ch/", "wss://relay.artio.inf.unibe.ch/"}
	miners := make([]*miner.RelayMiner, 0)
	for _, relay := range startingRelays {
		miners = append(miners, miner.NewMiner(relay))
	}

	neo := storage.Neo4jInstance{Username: "neo4j", Password: "fancyPassword", URI: "neo4j://localhost:7687", DBName: "neo4j"}
	err := neo.Init()
	if err != nil {
		log.Fatalf("Error on neo4j init: %v", err)
		return
	}
	_ = neo.Clean()

	defer neo.Close()
	manager := miner.Manager{Neo: &neo, MaxRecursion: 1}

	manager.Run(startingRelays)

}
