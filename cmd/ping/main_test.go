package main

import (
	"testing"
)

func TestHandler_Notifies(t *testing.T) {
	called := false

	ctx := Ctx{
		Notify: func(domain string) {
			called = true
		},
		GetUrls: func() []string {
			return []string{"https://google.com/404"}
		},
	}

	ctx.Run()

	if !called {
		t.Error("Did not notify available domain")
	}
}
