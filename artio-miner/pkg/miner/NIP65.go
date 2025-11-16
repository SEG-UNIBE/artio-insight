package miner

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"slices"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nbd-wtf/go-nostr"
)

/*
GetRelayList fetches all the Events of Type 10002 from the relay
*/
func GetRelayList(address string) ([]*nostr.Event, error) {
	interrupt := make(chan os.Signal, 1)
	eventList := make([]*nostr.Event, 0)
	signal.Notify(interrupt, os.Interrupt)
	c, _, err := websocket.DefaultDialer.Dial(address, nil)
	if err != nil {
		log.Println("dial:", err)
		return eventList, err
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			var response []json.RawMessage
			var messageType string
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read error:", err)
				return
			}
			_ = json.Unmarshal(message, &response)
			if len(response) > 0 {
				_ = json.Unmarshal(response[0], &messageType)
				if messageType == "EOSE" {
					// relay signaled EOSe
					// no more messages are coming -> closing session
					done <- struct{}{}
					return
				} else if messageType == "CLOSE" || messageType == "CLOSED" {
					// relay is closing the connection
					done <- struct{}{}
					return
				} else if messageType == "EVENT" {
					// unmarshall the event
					var event nostr.Event
					if err := json.Unmarshal(response[2], &event); err != nil {
						log.Printf("error while unamrshalling event: %s\n", err)

					}
					eventList = append(eventList, &event)
				} else if messageType == "NOTICE" {
					// relay sent a notice
					done <- struct{}{}
				} else if messageType == "AUTH" {
					// relay sent an auth message
					// we will ignore it for now
				} else {
					// we received a message that is not EOSE or EVENT, e.g. AUTH
					log.Printf("recv: %s", message)
				}
			}
		}
	}()

	err = c.WriteMessage(websocket.TextMessage, []byte("[\"REQ\", \"1\", {\"kinds\": [10002], \"limit\": 10000}]"))
	if err != nil {
		log.Println("write:", err)
	}

	go func() {
		time.Sleep(30 * time.Second)
		interrupt <- os.Signal(syscall.SIGINT)
	}()

	for {
		select {
		case <-done:
			return eventList, nil
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close error:", err)
				return eventList, err
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return eventList, nil
		}
	}
}

/*
FindNeighbours parses a list of nostr.Event to find all relays that are used by another user.
*/
func FindNeighbours(eventList []*nostr.Event) []string {
	wrongKind := 0
	var foundRelays []string
	// all the events must be of kind 10002
	for _, event := range eventList {
		if event.Kind != 10002 {
			wrongKind++
			continue
		}

		for _, tag := range event.Tags {
			if tagType := tag[0]; tagType == "r" {
				foundRelays = append(foundRelays, tag[1])
			}
		}
	}

	slices.Sort(foundRelays)
	return slices.Compact(foundRelays)
}
