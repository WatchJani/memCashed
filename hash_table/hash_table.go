package hash_table

import (
	"hash/fnv"
	"time"
)

type Engine struct {
	sendCh        []chan Payload
	shard         []map[string]Key
	numberOfShard uint32
}

type Key struct {
	field []byte
	ttl   time.Time
	//pointer on lru
}

type Payload struct {
	key   []byte
	field []byte
	ttl   uint32
}

func (e *Engine) Insert(key, field []byte, ttl uint32) {
	f := fnv.New32a()
	f.Write(key)

	shardIndex := int(f.Sum32() % uint32(e.numberOfShard))

	e.sendCh[shardIndex] <- Payload{
		key:   key,
		field: field,
		ttl:   ttl,
	}
}

func NewEngine(capacity uint32) *Engine {
	if capacity < 1 {
		capacity = 1
	}

	e := Engine{
		sendCh:        make([]chan Payload, capacity),
		shard:         make([]map[string]Key, capacity),
		numberOfShard: capacity,
	}

	for index := 0; index < int(capacity); index++ {
		e.sendCh[index] = make(chan Payload, 100) // Buffered channel
		e.shard[index] = make(map[string]Key)

		go e.ReceiveTask(index) //spawn new threads
	}

	return &e
}

func (e *Engine) ReceiveTask(index int) {
	for data := range e.sendCh[index] {
		
		e.shard[index][string(data.key)] = Key{}
	}
}
