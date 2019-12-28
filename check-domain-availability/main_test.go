package main

import (
	"os"
	"testing"
)

func TestHandler_Notifies(t *testing.T) {
	os.Setenv("WANTED_DOMAINS", "tzacwierjiyknoelkefbmyankdnlxbvaoujuizfy.com")

	called := false

	notify := func(domain string) {
		called = true
	}

	handler := CreateHandler(notify)

	handler()

	if !called {
		t.Error("Did not notify available domain")
	}
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
