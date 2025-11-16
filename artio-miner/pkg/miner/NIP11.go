package miner

import (
	"context"
	"io"
	"log"
	"net/http"
	"time"
)

/*
GetNip11 fetches the NIP 11 Information for a specifc relay
*/
func GetNip11(relay string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	client := &http.Client{Timeout: 3 * time.Second}
	method := "GET"

	req, err := http.NewRequestWithContext(ctx, method, relay, nil)
	if err != nil {
		log.Printf("Relay %s returned error: %s", relay, err)
		return nil, err
	}
	req.Header.Add("Accept", "application/nostr+json")
	// req.Header.Add("User-Agent", "relay-miner")
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Relay %s returned error: %s", relay, err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Relay %s returned error: %s", relay, err)
		return nil, err
	}
	return body, nil
}
