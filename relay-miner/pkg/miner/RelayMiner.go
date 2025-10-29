package miner

import (
	"log"

	"github.com/nbd-wtf/go-nostr"
)

/*
RelayList Object to store the response of a single relay
*/
type RelayMiner struct {
	Relay       string
	EventList   []*nostr.Event
	NIP11Result string
}

/*
LoadNIP11 Load the NIP-11 Result into the object
*/
func (rm *RelayMiner) LoadNIP11() {
	result, err := GetNip11(rm.Relay)
	if err != nil {
		log.Printf("error occured: %s/n", err)
		return
	}
	rm.NIP11Result = result
	return
}

/*
LoadRelayLists Load the NIP-11 Result into the object
*/
func (rm *RelayMiner) LoadRelayLists() {
	result, err := GetRelayList(rm.Relay)
	if err != nil {
		log.Printf("error occured: %s/n", err)
		return
	}
	rm.EventList = result
	return
}
