package main

import (
	"os"
	"testing"
)

func TestHandler(t *testing.T) {
	os.Setenv("WANTED_DOMAINS", "google.com,tzacwierjiyknoelkefbmyankdnlxbvaoujuizfy.com")

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
