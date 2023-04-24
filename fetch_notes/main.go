package main

import (
	"context"
	"fmt"
	"nostrgo/utils"
	"time"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
)

func main() {
	relay, err := nostr.RelayConnect(context.Background(), utils.DAMUS_RELAY)
	if err != nil {
		panic(err)
	}
	_, pubkeyHex, err := nip19.Decode(utils.KWINTEN_NPUB)
	if err != nil {
		panic(err)
	}
	since := time.Now().Add(time.Hour * 24 * 7 * -1)
	filters := []nostr.Filter{{
		Kinds:   []int{1},
		Authors: []string{pubkeyHex.(string)},
		Since:   &since,
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
