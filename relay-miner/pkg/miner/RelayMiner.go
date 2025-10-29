package miner

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip11"
)

/*
RelayMiner Object to store the response of a single relay
*/
type RelayMiner struct {
	Relay           string
	EventList       []*nostr.Event
	nip11Result     []byte // store both the raw result and the parsed to keep information that might not be compliant with NIP-11
	Nip11Document   *nip11.RelayInformationDocument
	NeighbourRelays []string
	loaded          bool
}

func (rm *RelayMiner) Load() {
	defer func() { rm.loaded = true }()
	rm.LoadNIP11()
	rm.LoadRelayLists()
	rm.LoadNeighbouringRelays()
}

/*
LoadNIP11 Load the NIP-11 Result into the object
*/
func (rm *RelayMiner) LoadNIP11() {
	address := fmt.Sprintf("https://%v/", rm.Relay)
	result, err := GetNip11(address)
	if err != nil {
		log.Printf("error occured: %s/n", err)
		return
	}
	rm.nip11Result = result
	rm.parseNip11()
	return
}

/*
parseNip11 internal method to parse the NIP11 document from string
*/
func (rm *RelayMiner) parseNip11() {
	byteNip11 := []byte(rm.nip11Result)
	var nipdoc nip11.RelayInformationDocument
	_ = json.Unmarshal(byteNip11, &nipdoc)
	rm.Nip11Document = &nipdoc

}

/*
LoadRelayLists Load the NIP-11 Result into the object
*/
func (rm *RelayMiner) LoadRelayLists() {
	address := fmt.Sprintf("wss://%v/", rm.Relay)
	result, err := GetRelayList(address)
	if err != nil {
		log.Printf("error occured: %s/n", err)
		return
	}
	rm.EventList = result
	return
}

func (rm *RelayMiner) LoadNeighbouringRelays() {
	neighbours := FindNeighbours(rm.EventList)
	rm.NeighbourRelays = neighbours
}

/*
Stats returns the basic stats of the relay to the command line
*/
func (rm *RelayMiner) Stats() {
	if !rm.loaded {
		rm.Load()
	}
	fmt.Printf("Relay: %v\n", rm.Relay)
	fmt.Printf("\tEvents: %v\n", len(rm.EventList))
	fmt.Printf("\tSoftare: %v\n", rm.Nip11Document.Software)
	fmt.Printf("\tNIPs: %v\n", rm.Nip11Document.SupportedNIPs)
	fmt.Printf("\tNeighbouring Relays: %v\n", len(rm.NeighbourRelays))
	//fmt.Printf("\tNeighbouring Relys: %v\n", rm.NeighbourRelays)

}

func NewMiner(relayUrl string) *RelayMiner {
	return &RelayMiner{Relay: relayUrl, EventList: make([]*nostr.Event, 0), loaded: false}
}
