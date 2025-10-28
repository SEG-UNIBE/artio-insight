package miner

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

/*
GetRelayList fetches all the Events of Type 10002 from the relay
*/
func GetRelayList(address string) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	c, _, err := websocket.DefaultDialer.Dial(address, nil)
	if err != nil {
		fmt.Println("dial:", err)
		return
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read error:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	err = c.WriteMessage(websocket.TextMessage, []byte("[\"COUNT\", \"1\", {}]"))
	if err != nil {
		log.Println("write:", err)
	}

	time.Sleep(5 * time.Second)
	interrupt <- os.Interrupt

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close error:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
