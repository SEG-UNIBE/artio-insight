package miner

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

/*
GetNip11 fetches the NIP 11 Information for a specifc relay
*/
func GetNip11(relay string) {
	client := &http.Client{}
	method := "GET"

	req, err := http.NewRequest(method, relay, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Accept", "application/nostr+json")
	// req.Header.Add("User-Agent", "relay-miner")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Println(string(body))
}
