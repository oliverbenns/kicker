package main

import (
	"context"
	"os"
	"testing"
)

func TestHandler(t *testing.T) {
	ctx := context.Background()
	os.Setenv("WANTED_DOMAINS", "google.com,tzacwierjiyknoelkefbmyankdnlxbvaoujuizfy.com")

	_, _ = Handler(ctx)
}

func TestIsDomainAvailable_Available(t *testing.T) {
	isAvailable := IsDomainAvailable("tzacwierjiyknoelkefbmyankdnlxbvaoujuizfy.com")

	if !isAvailable {
		t.Error("Incorrect evaluation of domain availability.")
	}
}

func TestIsDomainAvailable_Unavailable(t *testing.T) {
	isAvailable := IsDomainAvailable("google.com")

	if isAvailable {
		t.Error("Incorrect evaluation of domain availability.")
	}
}
