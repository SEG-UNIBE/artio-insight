package miner

import (
	"log"
	"time"

	"github.com/SEG-UNIBE/artio-insight/relay-miner/pkg/helper"
)

/*
Runner that will mine the RelayMiner objects
*/
type Runner struct {
	*Manager
	Id      int
	running bool
	idle    bool
}

/*
handleRelay process a single relay and store its information in the database
*/
func (rnr *Runner) handleRelay(relay *RelayMiner) {
	rnr.SetLoadMapEntryTrue(relay.CleanName())
	// load the relay information
	relay.Load()
	//relay.Stats()

	// merge the relay
	rnr.Neo.Execute(`MERGE(r:Relay {name: $name, isValid: $isValid, validReason: $validReason})`, map[string]any{"name": relay.CleanName(), "validReason": relay.InvalidReason, "isValid": relay.IsValid})
	rnr.Neo.Execute(`MERGE(r:RelayAlternativeName {name: $name})`, map[string]any{"name": relay.Relay})
	rnr.Neo.Execute(`MATCH(r:Relay), (ra:RelayAlternativeName) WHERE r.name=$name and ra.name=$alternativeName MERGE (r)-[:ALT_NAME]->(ra);`, map[string]any{"alternativeName": relay.Relay, "name": relay.CleanName()})
	if relay.DetectedBy != nil {
		rnr.Neo.Execute(`MATCH(r1:Relay), (r2:Relay) WHERE r1.name=$name1 and r2.name=$name2 MERGE (r1)-[:DETECTED]->(r2);`, map[string]any{"name1": relay.DetectedBy.CleanName(), "name2": relay.CleanName()})
	}
	if !relay.IsValid {
		return
	}

	// do the version
	rnr.Neo.Execute(`MERGE(s:Software {software: $software})`, map[string]any{"software": relay.Software()})

	// do the nip support
	if relay.Nip11Document != nil {
		for _, nip := range relay.Nip11Document.SupportedNIPs {
			rnr.Neo.Execute(`MATCH(r:Relay), (n:NIP) WHERE r.name=$name and n.name=$nip MERGE (r)-[:IMPLEMENTS]->(n);`, map[string]any{"nip": nip, "name": relay.CleanName()})
		}
	}

	// merge relation between relay and version
	rnr.Neo.Execute(`MATCH(r:Relay), (s:Software) WHERE r.name=$name and s.software=$version MERGE (r)-[:USES_SOFTWARE]->(s);`, map[string]any{"version": relay.Software(), "name": relay.CleanName()})

	// merge the public key of the owner
	rnr.Neo.Execute(`MERGE(u:User {pubkey: $pubkey})`, map[string]any{"pubkey": relay.PublicKey()})

	// merge relation between relay and owner
	rnr.Neo.Execute(`MATCH(r:Relay), (u:User) WHERE r.name=$name and u.pubkey=$pubkey MERGE (u)-[:OWNS]->(r);`, map[string]any{"pubkey": relay.PublicKey(), "name": relay.CleanName()})

	// do the IP addresses
	for _, ip := range relay.Ips {
		rnr.Neo.Execute(`MERGE(i:IP {address: $address})`, map[string]any{"address": ip.String()})
		rnr.Neo.Execute(`MATCH(r:Relay), (i:IP) WHERE r.name=$name and i.address=$address MERGE (r)-[:HAS_IP]->(i);`, map[string]any{"address": ip.String(), "name": relay.CleanName()})
	}

	if relay.RecursionLevel > 0 {
		log.Printf("Runner %d: Found %d new Relays for possible mining on %s\n", rnr.Id, len(relay.NeighbourRelays), relay.Relay)
		for _, rel := range relay.NeighbourRelays {
			// create the new RelayMiner object and enqueue it for further processing

			newRelay := NewMiner(rel)
			newRelay.DetectedBy = relay
			newRelay.RecursionLevel = relay.RecursionLevel - 1
			newRelay.Validate()
			rnr.Neo.Execute(`MERGE(r:Relay {name: $name, isValid: $isValid, validReason: $validReason})`, map[string]any{"name": newRelay.CleanName(), "validReason": newRelay.InvalidReason, "isValid": newRelay.IsValid})
			// rnr.Neo.Execute(`MERGE(r:Relay {name: $name})`, map[string]any{"name": newRelay.CleanName()})

			if rnr.GetLoadMapEntry(newRelay.CleanName()) {
				rnr.Neo.Execute(`MATCH(r1:Relay), (r2:Relay) WHERE r1.name=$name1 and r2.name=$name2 MERGE (r1)-[:DETECTED]->(r2);`, map[string]any{"name1": relay.CleanName(), "name2": newRelay.CleanName()})
			}

			rnr.Enqueue(newRelay)
		}
		if rnr.PushUsers {
			log.Printf("Runner %d: Found %d new NIP-65 messages\n", rnr.Id, len(relay.EventList))
			for _, evt := range relay.EventList {
				rnr.Neo.Execute(`MERGE(u:User {pubkey: $pubkey})`, map[string]any{"pubkey": evt.PubKey})

				pubkey, relays := helper.FindRelayForUser(evt)
				for _, rel := range relays {
					rnr.Neo.Execute(`MATCH(r:Relay), (u:User) WHERE r.name=$name and u.pubkey=$pubkey MERGE (u)-[:USES]->(r);`, map[string]any{"pubkey": pubkey, "name": helper.CleanRelayName(rel)})
				}
			}
		}
	}
}

func (rnr *Runner) Run() {
	rnr.running = true
	log.Printf("Runner %d started\n", rnr.Id)
	for rnr.running {
		nextMiner := rnr.Dequeue()
		if nextMiner == nil {
			if !rnr.idle {
				log.Printf("Runner %d is idle\n", rnr.Id)
			}
			rnr.idle = true

			time.Sleep(time.Second)
			continue
		} else {
			if rnr.idle {
				log.Printf("Runner %d is running\n", rnr.Id)
			}
			rnr.idle = false
			log.Printf("Runner %d is running with Relay %s\n", rnr.Id, nextMiner.Relay)
			rnr.handleRelay(nextMiner)
		}

	}
}

func (rnr *Runner) SignalEnd() {
	rnr.running = false
}

func (rnr *Runner) IsRunning() bool {
	return rnr.running
}

func (rnr *Runner) IsIdle() bool {
	return rnr.idle
}
