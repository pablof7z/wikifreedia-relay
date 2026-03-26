package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip77"
)

func TestRelayAdvertisesNegentropyInNIP11(t *testing.T) {
	relay, db, err := newRelay(filepath.Join(t.TempDir(), "relay.db"))
	if err != nil {
		t.Fatalf("failed to create relay: %v", err)
	}
	defer db.Close()

	server := httptest.NewServer(relay)
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Accept", "application/nostr+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to fetch nip-11: %v", err)
	}
	defer resp.Body.Close()

	var info struct {
		SupportedNIPs []any `json:"supported_nips"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		t.Fatalf("failed to decode nip-11: %v", err)
	}

	if !slices.ContainsFunc(info.SupportedNIPs, func(v any) bool {
		n, ok := v.(float64)
		return ok && int(n) == 77
	}) {
		t.Fatalf("expected supported_nips to contain 77, got %#v", info.SupportedNIPs)
	}
}

func TestRelayServesNegentropySync(t *testing.T) {
	relay, db, err := newRelay(filepath.Join(t.TempDir(), "relay.db"))
	if err != nil {
		t.Fatalf("failed to create relay: %v", err)
	}
	defer db.Close()

	sk := nostr.GeneratePrivateKey()
	evt := nostr.Event{
		CreatedAt: nostr.Now(),
		Kind:      30818,
		Tags:      nostr.Tags{{"d", "athens"}},
		Content:   "wiki article payload",
	}
	if err := evt.Sign(sk); err != nil {
		t.Fatalf("failed to sign event: %v", err)
	}

	if _, err := relay.AddEvent(context.Background(), &evt); err != nil {
		t.Fatalf("failed to add event: %v", err)
	}

	server := httptest.NewServer(relay)
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	idch, err := nip77.FetchIDsOnly(ctx, strings.Replace(server.URL, "http://", "ws://", 1), nostr.Filter{})
	if err != nil {
		t.Fatalf("failed to start negentropy sync: %v", err)
	}

	var ids []string
	for id := range idch {
		ids = append(ids, id)
	}

	if !slices.Contains(ids, evt.ID) {
		t.Fatalf("expected negentropy sync to return %s, got %#v", evt.ID, ids)
	}
}
