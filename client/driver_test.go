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

func Store(driver *Driver) error {
	_, err := driver.SetReq(key, value, ttl)
	if err != nil {
		return err
	}

	return nil
}

func TestGetReq(t *testing.T) {
	driver := New(":5000", 15)

	if err := driver.Init(); err != nil {
		t.Error(err)
	}

	if err := Store(driver); err != nil {
		t.Error(err)
	}

	res, err := driver.GetReq(key)
	if err != nil {
		t.Log(err)
	}

	get := string(<-res)
	expected := string(value)

	if get != expected {
		t.Errorf("expected: %s | get: %s", expected, get)
	}
}

func TestDeleteReq(t *testing.T) {
	driver := New(":5000", 15)

	if err := driver.Init(); err != nil {
		t.Error(err)
	}

	if err := Store(driver); err != nil {
		t.Error(err)
	}

	res, err := driver.DeleteReq(key)
	if err != nil {
		t.Log(err)
	}

	get := string(<-res)
	expected := "Deleted"

	if get != expected {
		t.Errorf("expected: %s | get: %s", expected, get)
	}

	//===========================================0

	res, err = driver.GetReq(key)
	if err != nil {
		t.Log(err)
	}

	get = string(<-res)
	expected = "object not found"

	if get != expected {
		t.Errorf("expected: %s | get: %s", expected, get)
	}
}
