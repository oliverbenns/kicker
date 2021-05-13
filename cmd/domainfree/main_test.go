package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestHandler_Notifies(t *testing.T) {
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "name,com,co\ntzacwierjiyknoelkefbmyankdnlxbvaoujuizfy,1,1")
		})
		http.ListenAndServe(":8080", nil)
	}()

	called := false

	notify := func(domain string) {
		called = true
	}

	handler := CreateHandler(notify, "http://localhost:8080")

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
