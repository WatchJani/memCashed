package hash_table

import (
	"errors"
	"hash/fnv"
	"log"
	"net"
	"root/link_list"
	"time"
	"unsafe"
)

var (
	ErrOperationIsNotSupported = errors.New("operation is not supported")
)

type Engine struct {
	sendCh        []chan interface{}
	shard         []map[string]Key
	numberOfShard uint32
}

type Key struct {
	field   []byte
	ttl     time.Time
	pointer *link_list.Node
}

// type Payload struct {
// 	key   []byte
// 	field []byte
// 	ttl   uint32
// }

func (e *Engine) Distribute(key []byte, setReq interface{}) {
	f := fnv.New32a()
	f.Write(key)

	shardIndex := int(f.Sum32() % uint32(e.numberOfShard))

	e.sendCh[shardIndex] <- setReq
	// e.sendCh[shardIndex] <- Payload{
	// 	key:   key,
	// 	field: field,
	// 	ttl:   ttl,
	// }
}

func NewEngine(capacity uint32) *Engine {
	if capacity < 1 {
		capacity = 1
	}

	e := Engine{
		sendCh:        make([]chan interface{}, capacity),
		shard:         make([]map[string]Key, capacity),
		numberOfShard: capacity,
	}

	for index := 0; index < int(capacity); index++ {
		e.sendCh[index] = make(chan interface{}, 100) // Buffered channel
		e.shard[index] = make(map[string]Key)

		go e.ReceiveTask(index) //spawn new threads
	}

	return &e
}

type SetReq struct {
	BaseReq
	field []byte
	ttl   uint32
}

func NewSetReq(
	key []byte,
	conn net.Conn,
	lru *link_list.DLL,
	field []byte,
	ttl uint32,
) SetReq {
	return SetReq{
		BaseReq: BaseReq{
			key:  key,
			conn: conn,
			lru:  lru,
		},
		field: field,
		ttl:   ttl,
	}
}

func NewDeleteReq(
	operation byte,
	key []byte,
	conn net.Conn,
	lru *link_list.DLL,
) BaseReq {
	return BaseReq{
		operation: operation,
		key:       key,
		conn:      conn,
		lru:       lru,
	}
}

func NewGetReq(
	operation byte,
	key []byte,
	conn net.Conn,
	lru *link_list.DLL,
) BaseReq {
	return BaseReq{
		operation: operation,
		key:       key,
		conn:      conn,
		lru:       lru,
	}
}

// For delete and get
type BaseReq struct {
	operation byte
	key       []byte
	conn      net.Conn
	lru       *link_list.DLL
}

type SysDelete struct {
	key []byte
}

func TLLParser(ttl uint32) time.Time {
	if ttl > 0 {
		return time.Now().Add(time.Duration(ttl) * time.Second)
	}

	return time.Time{}
}

func (e *Engine) ReceiveTask(index int) {
	for data := range e.sendCh[index] {
		switch v := data.(type) {
		case SetReq:
			obj := v

			key := string(obj.key)
			//insert in hash map
			e.shard[index][key] = Key{
				field: obj.field,
				ttl:   TLLParser(obj.ttl),
			}

			//insert in lru
			obj.lru.Inset(link_list.NewValue(unsafe.Pointer(&obj.field[0]), key))

			if _, err := obj.conn.Write([]byte("object inserted")); err != nil {
				log.Println(err)
			}
		case BaseReq:
			obj := v

			valueObject, isFound := e.shard[index][string(obj.key)]
			if !isFound {
				if _, err := obj.conn.Write([]byte("object not found")); err != nil {
					log.Println(err)
				}
				continue
			}

			if obj.operation == 'G' {
				if !valueObject.ttl.IsZero() && time.Now().After(valueObject.ttl) {
					delete(e.shard[index], string(v.key))
					obj.lru.Delete(valueObject.pointer) //delete from lru
					if _, err := obj.conn.Write([]byte("time expire")); err != nil {
						log.Println(err)
					}
				}

				if _, err := obj.conn.Write(valueObject.field); err != nil {
					log.Println(err)
				}
				continue
			}

			//for 'D' operation
			delete(e.shard[index], string(v.key))
			obj.lru.Delete(valueObject.pointer) //delete from lru
			if _, err := obj.conn.Write([]byte("Deleted")); err != nil {
				log.Println(err)
			}
		case SysDelete:
			delete(e.shard[index], string(v.key))
		default:
			log.Println(ErrOperationIsNotSupported)
		}
	}
}
