package helper

import (
	"net"
	"net/url"
	"strings"

	"github.com/nbd-wtf/go-nostr"
)

/*
CleanRelayName cleans a relay name by removing any ws://, wss:// and more
*/
func CleanRelayName(name string) string {
	name = strings.ReplaceAll(name, "ws://", "")
	name = strings.ReplaceAll(name, "wss://", "")
	name = strings.ReplaceAll(name, "http://", "")
	name = strings.ReplaceAll(name, "https://", "")
	name = strings.ReplaceAll(name, "/", "")
	name = strings.ReplaceAll(name, " ", "")
	name = strings.ReplaceAll(name, "\t", "")
	return name
}

/*
FindRelayForUser parses a list of nostr.Event to find all relays that are used by another user.
*/
func FindRelayForUser(event *nostr.Event) (string, []string) {
	var relays []string
	for _, tag := range event.Tags {
		if tagType := tag[0]; tagType == "r" {
			relays = append(relays, tag[1])
		}
	}
	return event.PubKey, relays
}

/*
ValidateURL validates if an url is valid and not only a localnetwork hostname
*/
func ValidateURL(uri string) (bool, string) {
	sharedIPNet := net.IPNet{
		IP:   net.IPv4(100, 64, 0, 0),
		Mask: net.CIDRMask(10, 32),
	}
	c, err := url.ParseRequestURI(uri)
	if err != nil {
		return false, "Invalid URL"
	}
	ipAddr := net.ParseIP(c.Hostname())
	if ipAddr.IsPrivate() {
		return false, "Private IP address"
	}
	if ipAddr.IsLoopback() {
		return false, "Loopback IP address"
	}
	if sharedIPNet.Contains(ipAddr) {
		return false, "Carrier-Grade NAT IP address"
	}
	if strings.HasSuffix(uri, ".onion") || strings.HasSuffix(uri, ".onion/") {
		return false, "TOR network address"
	}
	return true, ""
}

/*
ValidateDNS resolves a hostname to an IP address
*/
func ValidateDNS(hostname string) ([]net.IP, string) {
	var result []net.IP
	ips, err := net.LookupIP(hostname)
	if err != nil {
		return result, "DNS resolution failed"
	}
	for _, ip := range ips {
		result = append(result, ip)
	}
	return result, ""
}
