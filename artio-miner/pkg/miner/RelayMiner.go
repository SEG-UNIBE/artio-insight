package miner

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/SEG-UNIBE/artio-insight/relay-miner/pkg/helper"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip11"
)

/*
RelayMiner Object to store the response of a single relay
*/
type RelayMiner struct {
	Relay            string
	EventList        []*nostr.Event
	nip11Result      []byte // store both the raw result and the parsed to keep information that might not be compliant with NIP-11
	Nip11Document    *nip11.RelayInformationDocument
	NeighbourRelays  []string
	Ips              []net.IP
	DnsInValidReason string
	loaded           bool
	IsValid          bool
	InvalidReason    string
	DetectedBy       *RelayMiner
	RecursionLevel   int
}

func (rm *RelayMiner) Load() {
	defer func() { rm.loaded = true }()
	rm.IsValid, rm.InvalidReason = helper.ValidateURL(rm.Relay)
	if !rm.IsValid {
		log.Println(rm.InvalidReason, ": ", rm.Relay)
		return
	}
	rm.Ips, rm.DnsInValidReason = helper.ValidateDNS(rm.CleanName())
	if rm.DnsInValidReason != "" {
		rm.IsValid = false
		rm.InvalidReason = rm.DnsInValidReason
		log.Println(rm.DnsInValidReason, ": ", rm.Relay)
		return
	}

	rm.LoadNIP11()
	if rm.RecursionLevel > 0 {
		rm.LoadRelayLists()
		rm.LoadNeighbouringRelays()
	}
}

func (rm *RelayMiner) Validate() {
	rm.IsValid, rm.InvalidReason = helper.ValidateURL(rm.Relay)
	if !rm.IsValid {
		log.Println(rm.InvalidReason, ": ", rm.Relay)
		return
	}
	rm.Ips, rm.DnsInValidReason = helper.ValidateDNS(rm.CleanName())
	if rm.DnsInValidReason != "" {
		rm.IsValid = false
		rm.InvalidReason = rm.DnsInValidReason
		log.Println(rm.DnsInValidReason, ": ", rm.Relay)
		return
	}
	return
}

/*
LoadNIP11 Load the NIP-11 Result into the object
*/
func (rm *RelayMiner) LoadNIP11() {
	var address string
	if strings.HasPrefix(rm.Relay, "ws://") {
		address = fmt.Sprintf("http://%v/", rm.CleanName())
	} else {
		address = fmt.Sprintf("https://%v/", rm.CleanName())
	}
	result, err := GetNip11(address)
	if err != nil {
		log.Printf("error occured: %s\n", err)
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
	address := fmt.Sprintf("%v", rm.Relay)
	result, err := GetRelayList(address)
	if err != nil {
		log.Printf("error occured: %s\n", err)
		return
	}
	rm.EventList = result
	return
}

/*
GetCleanRelayList Get cleaned up list of relays to load the NIP-11 from.
*/
func (rm *RelayMiner) GetCleanRelayList() []string {
	outputList := make([]string, 0)
	for _, relay := range rm.NeighbourRelays {
		relay = helper.CleanRelayName(relay)
		outputList = append(outputList, relay)
	}
	return outputList
}

/*
Software gets the software stack from the NIP 11 document
*/
func (rm *RelayMiner) Software() string {
	if rm.Nip11Document == nil {
		return "N/A"
	}
	return rm.Nip11Document.Software
}

/*
CleanName returns the cleaned name of the relay
*/
func (rm *RelayMiner) CleanName() string {
	return helper.CleanRelayName(rm.Relay)
}

/*
PublicKey returns the public key of the relay if available
*/
func (rm *RelayMiner) PublicKey() string {
	if rm.Nip11Document == nil {
		return "N/A"
	}
	return rm.Nip11Document.PubKey
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
	if rm.Nip11Document != nil {
		fmt.Printf("\tSoftare: %v\n", rm.Nip11Document.Software)
		fmt.Printf("\tNIPs: %v\n", rm.Nip11Document.SupportedNIPs)
	} else {
		fmt.Printf("\tSoftare: %v\n", "N/A")
		fmt.Printf("\tNIPs: %v\n", "N/A")
	}
	fmt.Printf("\tNeighbouring Relays: %v\n", len(rm.NeighbourRelays))
	//fmt.Printf("\tNeighbouring Relys: %v\n", rm.NeighbourRelays)

}

func NewMiner(relayUrl string) *RelayMiner {
	return &RelayMiner{Relay: relayUrl, EventList: make([]*nostr.Event, 0), loaded: false, DetectedBy: nil}
}
