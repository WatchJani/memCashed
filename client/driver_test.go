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

	_, err := driver.SetReq(key, value, ttl)
	if err != nil {
		t.Error(err)
	}
}
