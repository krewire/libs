package storage

import "testing"

func TestLocalRoundTripAndTraversalProtection(t *testing.T) {
	store, err := NewLocal(t.TempDir(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Put("objects/item.txt", []byte("payload")); err != nil {
		t.Fatal(err)
	}
	got, err := store.Get("objects/item.txt")
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "payload" {
		t.Fatalf("payload = %q, want payload", got)
	}
	info, err := store.Stat("objects/item.txt")
	if err != nil {
		t.Fatal(err)
	}
	if info.Key != "objects/item.txt" || info.Size != int64(len("payload")) {
		t.Fatalf("info = %+v", info)
	}
	if err := store.Put("../escape.txt", []byte("blocked")); err == nil {
		t.Fatal("path traversal must be rejected")
	}
	if err := store.Delete("objects/item.txt"); err != nil {
		t.Fatal(err)
	}
}
