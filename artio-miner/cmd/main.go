package main

import (
	"log"
	"os"

	"github.com/SEG-UNIBE/artio-insight/relay-miner/pkg/miner"
	"github.com/SEG-UNIBE/artio-insight/relay-miner/pkg/storage"
	"github.com/joho/godotenv"
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
	err := godotenv.Load(".env")

	uri := os.Getenv("NEO4J_URI")
	password := os.Getenv("NEO4J_PASSWORD")
	db := os.Getenv("NEO4J_DB")
	username := os.Getenv("NEO4J_USERNAME")

	neo := storage.Neo4jInstance{Username: username, Password: password, URI: uri, DBName: db}
	err := neo.Init()
	if err != nil {
		log.Fatalf("Error on neo4j init: %v", err)
		return
	}
	_ = neo.Clean()

	defer neo.Close()
	manager := miner.Manager{Neo: &neo, MaxRecursion: 2, MaxRunners: 2, PushUsers: false}

	manager.Run(startingRelays)

}
