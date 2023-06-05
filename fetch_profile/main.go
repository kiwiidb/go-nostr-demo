package main

import (
	"context"
	"fmt"
	"nostrgo/utils"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
)

func main() {
	relay, err := nostr.RelayConnect(context.Background(), utils.DAMUS_RELAY)
	if err != nil {
		panic(err)
	}
	_, pubkeyHex, err := nip19.Decode(utils.KWINTEN_NPUB)
	fmt.Println(pubkeyHex)
	if err != nil {
		panic(err)
	}
	filters := []nostr.Filter{{
		Kinds:   []int{0},
		Authors: []string{pubkeyHex.(string)},
	}}
	fmt.Println("Starting subscription")
	sub := relay.Subscribe(context.Background(), filters)
	for {
		select {
		case ev := <-sub.Events:
			fmt.Println(ev.Content)
		case <-sub.EndOfStoredEvents:
			return
		}
	}
}
