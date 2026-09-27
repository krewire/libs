package storage

import (
	"context"
	"errors"
	"testing"
)

func TestMemoryContextStore(t *testing.T) {
	store, err := NewMemory()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	if err := store.PutContext(ctx, "item", []byte("data")); err != nil {
		t.Fatal(err)
	}
	got, err := store.GetContext(ctx, "item")
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "data" {
		t.Fatalf("GetContext = %q", got)
	}
}

func TestCanceledContextDoesNotWrite(t *testing.T) {
	store, err := NewMemory()
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := store.PutContext(ctx, "item", []byte("data")); !errors.Is(err, context.Canceled) {
		t.Fatalf("PutContext error = %v, want context.Canceled", err)
	}
}
