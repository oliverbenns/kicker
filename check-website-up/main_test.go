package main

import (
	"os"
	"testing"
)

func TestHandler_Notifies(t *testing.T) {
	os.Setenv("CONCERNED_ENDPOINTS", "https://google.com/404")

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
