package helper

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/nbd-wtf/go-nostr"
)

/*
TestCleanRelayName tests the CleanRelayName function and its output
*/
func TestCleanRelayName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "CleanRelayName1", args: args{"wss://relay.relay.com/"}, want: "relay.relay.com"},
		{name: "CleanRelayName2", args: args{"ws://relay.relay.com/"}, want: "relay.relay.com"},
		{name: "CleanRelayName3", args: args{"wss://relay.relay.com"}, want: "relay.relay.com"},
		{name: "CleanRelayName3", args: args{"ws://relay.relay.com"}, want: "relay.relay.com"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CleanRelayName(tt.args.name); got != tt.want {
				t.Errorf("CleanRelayName() = %v, want %v", got, tt.want)
			}
		})
	}
}

/*
TestFindRelayForUser tests the FindRelayForUser function and its output
*/
func TestFindRelayForUser(t *testing.T) {
	type args struct {
		event *nostr.Event
	}
	emptyEvent := &nostr.Event{}
	populatedEvent := &nostr.Event{
		PubKey: "testpubkey",
		Tags: nostr.Tags{
			nostr.Tag{"r", "wss://relay1.com/"},
			nostr.Tag{"r", "wss://relay2.com/"},
		},
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 []string
	}{
		{name: "FindRelayForUser_EmptyEvent", args: args{event: emptyEvent}, want: "", want1: []string{}},
		{name: "FindRelayForUser_FullEvent", args: args{event: populatedEvent}, want: "testpubkey", want1: []string{"wss://relay1.com/", "wss://relay2.com/"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := FindRelayForUser(tt.args.event)
			if got != tt.want {
				t.Errorf("FindRelayForUser() pubkey got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(len(got1), len(tt.want1)) {
				t.Errorf("FindRelayForUser() list got = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestValidateURL(t *testing.T) {
	type args struct {
		uri string
	}
	tests := []struct {
		name       string
		args       args
		wantBool   bool
		wantString string
	}{
		{name: "ValidateURL1", args: args{uri: "relay.relay.com"}, wantBool: true, wantString: ""},
		{name: "ValidateURL2", args: args{uri: "relay.relay.com"}, wantBool: true, wantString: ""},
		{name: "ValidateURL3", args: args{uri: "relay.relay.com"}, wantBool: true, wantString: ""},
		{name: "ValidateURL4", args: args{uri: "relay.relay.com"}, wantBool: true, wantString: ""},
		{name: "ValidateURL5", args: args{uri: "192.168.10.10"}, wantBool: false, wantString: "Private IP address"},
		{name: "ValidateURL6", args: args{uri: "10.10.10.10"}, wantBool: false, wantString: "Private IP address"},
		{name: "ValidateURL7", args: args{uri: "100.64.224.5"}, wantBool: false, wantString: "Carrier-Grade NAT IP address"},
		{name: "ValidateURL8", args: args{uri: "127.0.0.1"}, wantBool: false, wantString: "Loopback IP address"},
		{name: "ValidateURL9", args: args{uri: "ex3znuu3kt4se7fjhc2l7zbjv2ydsajqi5suegk3gpfuqlzdgtl4f3qd.onion"}, wantBool: false, wantString: "TOR network address"},
	}
	for _, tt := range tests {
		for _, prefix := range []string{"ws://", "wss://", "http://", "https://"} {
			t.Run(tt.name, func(t *testing.T) {
				value := fmt.Sprintf("%s%s/", prefix, tt.args.uri)
				if gotBool, gotString := ValidateURL(value); gotBool != tt.wantBool || (gotString != tt.wantString) {
					t.Errorf("ValidateURL(%v) = %v with %v, want %v with %v", value, gotBool, gotString, tt.wantBool, tt.wantString)
				}
			})
		}
	}
}
