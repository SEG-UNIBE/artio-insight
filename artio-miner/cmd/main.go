package main

import (
	"fmt"
	"log"

	"github.com/SEG-UNIBE/artio-insight/relay-miner/pkg/helper"
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

	for _, relay := range miners {
		handleRelay(relay, neo, true)
	}
}

func handleRelay(relay *miner.RelayMiner, neo storage.Neo4jInstance, recursive bool) {
	if !helper.ValidateURL(relay.Relay) {
		fmt.Println("invalid relay url:", relay.Relay)
		return
	}
	relay.Load()
	relay.Stats()

	// merge the relay
	neo.Execute(`MERGE(r:Relay {name: $name})`, map[string]any{"name": relay.CleanName()})
	neo.Execute(`MERGE(r:RelayAlternativeName {name: $name})`, map[string]any{"name": relay.Relay})
	neo.Execute(`MATCH(r:Relay), (ra:RelayAlternativeName) WHERE r.name=$name and ra.name=$alternativeName MERGE (r)-[:ALT_NAME]->(ra);`, map[string]any{"alternativeName": relay.Relay, "name": relay.CleanName()})

	// do the version
	neo.Execute(`MERGE(s:Software {software: $software})`, map[string]any{"software": relay.Software()})

	// merge relation between relay and version
	neo.Execute(`MATCH(r:Relay), (s:Software) WHERE r.name=$name and s.software=$version MERGE (r)-[:USES_SOFTWARE]->(s);`, map[string]any{"version": relay.Software(), "name": relay.CleanName()})

	// merge the public key of the owner
	neo.Execute(`MERGE(u:User {pubkey: $pubkey})`, map[string]any{"pubkey": relay.PublicKey()})

	// merge relation between relay and owner
	neo.Execute(`MATCH(r:Relay), (u:User) WHERE r.name=$name and u.pubkey=$pubkey MERGE (u)-[:OWNS]->(r);`, map[string]any{"pubkey": relay.PublicKey(), "name": relay.CleanName()})

	if recursive {
		for _, rel := range relay.NeighbourRelays {
			neo.Execute(`MERGE(r:Relay {name: $name})`, map[string]any{"name": rel})
			relayMiner := miner.NewMiner(rel)
			go handleRelay(relayMiner, neo, false)
		}

		for _, evt := range relay.EventList {
			neo.Execute(`MERGE(u:User {pubkey: $pubkey})`, map[string]any{"pubkey": evt.PubKey})

			pubkey, relays := helper.FindRelayForUser(evt)
			for _, rel := range relays {
				neo.Execute(`MATCH(r:Relay), (u:User) WHERE r.name=$name and u.pubkey=$pubkey MERGE (u)-[:USES]->(r);`, map[string]any{"pubkey": pubkey, "name": rel})
			}
		}
	}
}
