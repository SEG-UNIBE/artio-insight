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
	startingRelays := []string{"relay.artiostr.ch", "relay.artio.inf.unibe.ch"}
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

	for _, relay := range miners {
		relay.Load()
		relay.Stats()

		for _, rel := range relay.NeighbourRelays {
			neo.Execute(`MERGE(r:Relay {name: $name})`, map[string]any{"name": rel})
		}

		for _, evt := range relay.EventList {
			neo.Execute(`MERGE(u:User {pubkey: $pubkey})`, map[string]any{"pubkey": evt.PubKey})

			pubkey, relays := miner.FindRelayForUser(evt)
			for _, rel := range relays {
				neo.Execute(`MATCH(r:Relay), (u:User) WHERE r.name=$name and u.pubkey=$pubkey MERGE (u)-[:USES]->(r);`, map[string]any{"pubkey": pubkey, "name": rel})
			}
		}
	}
}
