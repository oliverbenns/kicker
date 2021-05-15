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

	ctx := Ctx{
		Notify: func(domain string) {
			called = true
		},
		GetCsvUrl: func() string {
			return "http://localhost:8080"
		},
	}

	ctx.Run()

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
