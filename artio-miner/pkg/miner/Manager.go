package miner

import (
	"github.com/SEG-UNIBE/artio-insight/relay-miner/pkg/helper"
	"github.com/SEG-UNIBE/artio-insight/relay-miner/pkg/storage"
)

/*
Manager is the main object to handle the mining process
*/
type Manager struct {
	Neo          *storage.Neo4jInstance
	MaxRecursion int
	miners       []*RelayMiner
	loadMap      map[string]bool
}

/*
Run starts the mining process for the given relays
*/
func (mgmt *Manager) Run(relays []string) {
	mgmt.loadMap = make(map[string]bool)
	for _, relay := range relays {
		mgmt.miners = append(mgmt.miners, NewMiner(relay))
	}

	for _, relay := range mgmt.miners {
		mgmt.handleRelay(relay, nil, mgmt.MaxRecursion)
	}
}

/*
handleRelay process a single relay and store its information in the database
*/
func (mgmt *Manager) handleRelay(relay *RelayMiner, detectedBy *RelayMiner, recursionLevel int) {
	if mgmt.loadMap[relay.CleanName()] {
		return
	}
	mgmt.loadMap[relay.CleanName()] = true

	// load the relay information
	relay.Load()
	relay.Stats()

	// merge the relay
	mgmt.Neo.Execute(`MERGE(r:Relay {name: $name, isValid: $isValid, validReason: $validReason})`, map[string]any{"name": relay.CleanName(), "validReason": relay.InvalidReason, "isValid": relay.IsValid})
	mgmt.Neo.Execute(`MERGE(r:RelayAlternativeName {name: $name})`, map[string]any{"name": relay.Relay})
	mgmt.Neo.Execute(`MATCH(r:Relay), (ra:RelayAlternativeName) WHERE r.name=$name and ra.name=$alternativeName MERGE (r)-[:ALT_NAME]->(ra);`, map[string]any{"alternativeName": relay.Relay, "name": relay.CleanName()})
	if !relay.IsValid {
		return
	}

	// do the version
	mgmt.Neo.Execute(`MERGE(s:Software {software: $software})`, map[string]any{"software": relay.Software()})

	// merge relation between relay and version
	mgmt.Neo.Execute(`MATCH(r:Relay), (s:Software) WHERE r.name=$name and s.software=$version MERGE (r)-[:USES_SOFTWARE]->(s);`, map[string]any{"version": relay.Software(), "name": relay.CleanName()})

	// merge the public key of the owner
	mgmt.Neo.Execute(`MERGE(u:User {pubkey: $pubkey})`, map[string]any{"pubkey": relay.PublicKey()})

	// merge relation between relay and owner
	mgmt.Neo.Execute(`MATCH(r:Relay), (u:User) WHERE r.name=$name and u.pubkey=$pubkey MERGE (u)-[:OWNS]->(r);`, map[string]any{"pubkey": relay.PublicKey(), "name": relay.CleanName()})

	// if we have a detectedBy, create the relationship
	if detectedBy != nil {
		mgmt.Neo.Execute(`MATCH(r1:Relay), (r2:Relay) WHERE r1.name=$name1 and r2.name=$name2 MERGE (r1)-[:DETECTED]->(r2);`, map[string]any{"name1": detectedBy.CleanName(), "name2": relay.CleanName()})
	}

	if recursionLevel > 0 {
		recursionLevel--
		for _, rel := range relay.NeighbourRelays {
			mgmt.Neo.Execute(`MERGE(r:Relay {name: $name})`, map[string]any{"name": rel})
			relayMiner := NewMiner(rel)

			go mgmt.handleRelay(relayMiner, relay, recursionLevel)
		}

		for _, evt := range relay.EventList {
			mgmt.Neo.Execute(`MERGE(u:User {pubkey: $pubkey})`, map[string]any{"pubkey": evt.PubKey})

			pubkey, relays := helper.FindRelayForUser(evt)
			for _, rel := range relays {
				mgmt.Neo.Execute(`MATCH(r:Relay), (u:User) WHERE r.name=$name and u.pubkey=$pubkey MERGE (u)-[:USES]->(r);`, map[string]any{"pubkey": pubkey, "name": rel})
			}
		}
	}
}
