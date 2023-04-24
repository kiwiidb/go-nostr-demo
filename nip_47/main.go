package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"nostrgo/utils"
	"os"
	"time"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip04"
)

func main() {

	nip47Privkey := os.Getenv("NIP_47_PRIV")
	nip47RecipientPubkey := os.Getenv("NIP_47_RECIPIENT_PUB")
	pub, err := nostr.GetPublicKey(nip47Privkey)
	if err != nil {
		panic(err)
	}

	ss, err := nip04.ComputeSharedSecret(nip47RecipientPubkey, nip47Privkey)
	if err != nil {
		panic(err)
	}
	secretMessage, err := ioutil.ReadFile("nip47_request.json")
	if err != nil {
		panic(err)
	}
	encryptedContent, err := nip04.Encrypt(string(secretMessage), ss)
	if err != nil {
		panic(err)
	}
	ev := nostr.Event{
		PubKey:    pub,
		CreatedAt: time.Now(),
		Kind:      23194,
		Tags:      nostr.Tags{nostr.Tag{"p", nip47RecipientPubkey}},
		Content:   encryptedContent,
	}

	// calling Sign sets the event ID field and the event Sig field
	ev.Sign(nip47Privkey)

	// publish the event to multiple relays
	for _, url := range []string{utils.NIP47_RELAY} {
		relay, err := nostr.RelayConnect(context.Background(), url)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("published to ", url, relay.Publish(context.Background(), ev))
	}
}
