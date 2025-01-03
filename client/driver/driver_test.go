package client

// import (
// 	"bytes"
// 	"log"
// 	"testing"

// 	p "github.com/WatchJani/memCashed/client/parser"
// )

// var (
// 	key   = []byte("key")
// 	value = []byte("value")
// 	ttl   = -1
// )

// func TestSetReq(t *testing.T) {
// 	driver, err := NewConnection(":5000", 15)

// 	if err != nil {
// 		t.Error(err)
// 	}

// 	res, err := driver.SetReq(key, value, ttl)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	expect := "object inserted"
// 	get := string(<-res)

// 	if expect != get {
// 		t.Errorf("expected: %s | get: %s", expect, get)
// 	}
// }

// func Store(driver *Connection) error {
// 	_, err := driver.SetReq(key, value, ttl)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func TestGetReq(t *testing.T) {
// 	driver, err := NewConnection(":5000", 15)

// 	if err != nil {
// 		t.Error(err)
// 	}

// 	if err := Store(driver); err != nil {
// 		t.Error(err)
// 	}

// 	res, err := driver.GetReq(key)
// 	if err != nil {
// 		t.Log(err)
// 	}

// 	get := string(<-res)
// 	expected := string(value)

// 	if get != expected {
// 		t.Errorf("expected: %s | get: %s", expected, get)
// 	}
// }

// func TestDeleteReq(t *testing.T) {
// 	driver, err := NewConnection(":5000", 15)

// 	if err != nil {
// 		t.Error(err)
// 	}

// 	if err := Store(driver); err != nil {
// 		t.Error(err)
// 	}

// 	res, err := driver.DeleteReq(key)
// 	if err != nil {
// 		t.Log(err)
// 	}

// 	get := string(<-res)
// 	expected := "deleted"

// 	if get != expected {
// 		t.Errorf("expected: %s | get: %s", expected, get)
// 	}

// 	//===========================================0

// 	res, err = driver.GetReq(key)
// 	if err != nil {
// 		t.Log(err)
// 	}

// 	get = string(<-res)
// 	expected = "object not found"

// 	if get != expected {
// 		t.Errorf("expected: %s | get: %s", expected, get)
// 	}
// }

// func TestEncode(t *testing.T) {
// 	get, err := p.Encode('S', key, value, ttl)
// 	if err != nil {
// 		t.Fail()
// 	}

// 	want := []byte{18, 0, 0, 0, 83, 3, 255, 255, 255, 255, 5, 0, 0, 0, 107, 101, 121, 118, 97, 108, 117, 101}

// 	if !bytes.Equal(get, want) {
// 		t.Errorf("wanted %b | get %b", want, get)
// 	}
// }

// func BenchmarkEncode(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		p.Encode('S', key, value, ttl)
// 	}
// }

// func BenchmarkGetReq(b *testing.B) {
// 	b.StopTimer()
// 	driver, err := NewConnection(":5000", 15)

// 	if err != nil {
// 		b.Error(err)
// 		return
// 	}

// 	b.StartTimer()

// 	for i := 0; i < b.N; i++ {
// 		go func() {
// 			res, err := driver.SetReq(key, value, ttl)
// 			if err != nil {
// 				log.Println(err)
// 			}

// 			<-res
// 		}()
// 	}
// }
