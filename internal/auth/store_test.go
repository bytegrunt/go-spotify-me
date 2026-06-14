package auth

import (
	"testing"
)

func TestNewOSTokenStore(t *testing.T) {
	store := NewOSTokenStore()
	if store == nil {
		t.Fatalf("Expected non-nil TokenStore")
	}
	
	_, ok := store.(*osTokenStore)
	if !ok {
		t.Errorf("Expected store to be of type *osTokenStore")
	}
}
