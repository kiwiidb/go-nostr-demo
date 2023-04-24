package main

import (
	"context"
	"fmt"
	"nostrgo/utils"
	"time"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip04"
	"github.com/nbd-wtf/go-nostr/nip19"
)

func main() {

	pub, err := nostr.GetPublicKey(utils.MY_PRIVKEY)
	if err != nil {
		panic(err)
	}
	fmt.Printf("My npub is %s \n", pub)

	_, recipientPubkeyHex, err := nip19.Decode(utils.KWINTEN_NPUB)
	if err != nil {
		panic(err)
	}

	ss, err := nip04.ComputeSharedSecret(recipientPubkeyHex.(string), utils.MY_PRIVKEY)
	if err != nil {
		panic(err)
	}
	secretMessage := "hello from antwerp: encrypted"
	encryptedContent, err := nip04.Encrypt(secretMessage, ss)
	if err != nil {
		panic(err)
	}
	ev := nostr.Event{
		PubKey:    pub,
		CreatedAt: time.Now(),
		Kind:      4,
		Tags:      nostr.Tags{nostr.Tag{"p", recipientPubkeyHex.(string)}},
		Content:   encryptedContent,
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
