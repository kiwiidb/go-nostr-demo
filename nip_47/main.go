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
	var filters nostr.Filters
	filters = []nostr.Filter{{
		Kinds:   []int{13194, 23195},
		Authors: []string{nip47RecipientPubkey},
	}}
	relay, err := nostr.RelayConnect(context.Background(), utils.NIP47_RELAY)
	if err != nil {
		panic(err)
	}
	sub := relay.Subscribe(context.Background(), filters)

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

	status := relay.Publish(context.Background(), ev)
	fmt.Println(status)
	for ev := range sub.Events {
		if ev.Kind == 23195 {
			decrypted, err := nip04.Decrypt(ev.Content, ss)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(string(decrypted))
		}
	}
}
