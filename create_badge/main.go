package main

import (
	"context"
	"fmt"
	"nostrgo/utils"
	"os"
	"time"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
)

func main() {

	_, privkey, err := nip19.Decode(os.Getenv("NSEC"))
	if err != nil {
		panic(err)
	}
	pub, err := nostr.GetPublicKey(privkey.(string))
	if err != nil {
		panic(err)
	}
	fmt.Printf("My pubkey is %s \n", pub)
	ev := nostr.Event{
		PubKey:    pub,
		CreatedAt: time.Now(),
		Kind:      30009,
		Tags: nostr.Tags{
			nostr.Tag{"d", "alby_btc_prague"},
			nostr.Tag{"name", "Alby @ BTC Prague 2023"},
			nostr.Tag{"image", "https://cdn.getalby-assets.com/nostr-badges/btc_prague_2048x2048.png", "2048x2048"},
			nostr.Tag{"image", "https://cdn.getalby-assets.com/nostr-badges/btc_prague_256x256.png", "256x256"},
		},
	}

	// calling Sign sets the event ID field and the event Sig field
	ev.Sign(privkey.(string))

	// publish the event to multiple relays
	for _, url := range []string{utils.DAMUS_RELAY, utils.BITCOINER_SOCIAL_RELAY, utils.NOS_LOL_RELAY} {
		relay, err := nostr.RelayConnect(context.Background(), url)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("published to ", url, relay.Publish(context.Background(), ev))
	}
}
