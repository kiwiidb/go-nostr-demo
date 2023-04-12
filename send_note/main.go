package main

import (
	"context"
	"fmt"
	"nostrgo/utils"
	"time"

	"github.com/nbd-wtf/go-nostr"
)

func main() {

	pub, err := nostr.GetPublicKey(utils.MY_PRIVKEY)
	if err != nil {
		panic(err)
	}
	fmt.Printf("My npub is %s \n", pub)
	ev := nostr.Event{
		PubKey:    pub,
		CreatedAt: time.Now(),
		Kind:      1,
		Tags:      nil,
		Content:   "TODO: fill in",
	}

	// calling Sign sets the event ID field and the event Sig field
	ev.Sign(utils.MY_PRIVKEY)

	// publish the event to multiple relays
	for _, url := range []string{utils.DAMUS_RELAY, utils.WIZ_RELAY, utils.EDEN_RELAY} {
		relay, err := nostr.RelayConnect(context.Background(), url)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("published to ", url, relay.Publish(context.Background(), ev))
	}
}
