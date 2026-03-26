package main

import (
	"context"
	"fmt"
	"net/http"
	"slices"

	"github.com/fiatjaf/eventstore/sqlite3"
	"github.com/fiatjaf/khatru"
	"github.com/nbd-wtf/go-nostr"
)

var allowedKinds = [...]int{30818, 30819, 818, 819}

func newRelay(databaseURL string) (*khatru.Relay, *sqlite3.SQLite3Backend, error) {
	relay := khatru.NewRelay()
	relay.Negentropy = true

	db := &sqlite3.SQLite3Backend{DatabaseURL: databaseURL}
	if err := db.Init(); err != nil {
		return nil, nil, err
	}

	relay.Info.Name = "Wikifreedia relay"
	relay.Info.PubKey = "fa984bd7dbb282f07e16e7ae87b26a2a7b9b90b7246a44771f0cf5ae58018f52"
	relay.Info.Description = "This is a relay for wiki events. It supports search."
	relay.Info.Icon = "https://cdn.satellite.earth/064536fe832f87eb16e113ff227b61eb34ae4cd8e4ece3ef2b67d71257e52c71.png"

	relay.StoreEvent = append(relay.StoreEvent, db.SaveEvent)
	relay.QueryEvents = append(relay.QueryEvents, db.QueryEvents)
	relay.CountEvents = append(relay.CountEvents, db.CountEvents)
	relay.DeleteEvent = append(relay.DeleteEvent, db.DeleteEvent)
	relay.ReplaceEvent = append(relay.ReplaceEvent, db.ReplaceEvent)
	relay.RejectEvent = append(relay.RejectEvent,
		func(ctx context.Context, event *nostr.Event) (reject bool, msg string) {
			if !slices.Contains(allowedKinds[:], event.Kind) {
				return true, "only wiki events are allowed here."
			}
			return false, "" // anyone else can
		},
	)

	return relay, db, nil
}

func main() {
	relay, _, err := newRelay("./db")
	if err != nil {
		panic(err)
	}

	fmt.Println("running on :3334")
	if err := http.ListenAndServe(":3334", relay); err != nil {
		panic(err)
	}
}
