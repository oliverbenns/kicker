package main

import (
	"testing"
)

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
