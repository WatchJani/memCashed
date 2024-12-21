package client

import (
	"testing"
)

var (
	key   = []byte("key")
	value = []byte("value")
	ttl   = -1
)

func TestSetReq(t *testing.T) {
	driver := New(":5000", 15)

	if err := driver.Init(); err != nil {
		t.Error(err)
	}

	res, err := driver.SetReq(key, value, ttl)
	if err != nil {
		t.Error(err)
	}

	expect := "object inserted"
	get := string(<-res)

	if expect != get {
		t.Errorf("expected: %s | get: %s", expect, get)
	}
}
