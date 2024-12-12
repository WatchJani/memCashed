package hash_table

import (
	"fmt"
	"hash/fnv"
	"time"
)

func main() {
	m := make(map[string]int)

	var (
		stop  bool = true
		index int
	)

	s := make([]string, 150_000_000)
	for index := range len(s) {
		s[index] = fmt.Sprintf("%d", index)
	}

	go func() {
		<-time.Tick(30 * time.Second)

		stop = false
	}()

	// 21s da mi izracuna ovo sve

	for stop {
		m[s[index%150_000_000]] = index
		// c.Insert(s[index%150_000_000], index)
		// m.Set(fmt.Sprintf("%d", index), index)
		index++
	}

	fmt.Println(index)
}

// start := time.Now()

// 	var key map[int]int = map[int]int{}
// 	hash := fnv.New32a()

// 	for index := 0; index < 150_000_000; index++ {
// 		hash.Write([]byte(fmt.Sprintf("%d", index)))
// 		shardIndex := int(hash.Sum32() % uint32(12))

// 		value, ok := key[shardIndex]

// 		if ok {
// 			value++
// 		}

// 		key[shardIndex] = value
// 	}

// 	fmt.Println(time.Since(start))

// 	for key, value := range key {
// 		percentage := float64(value) / float64(150_000_000) * 100
// 		fmt.Printf("Ključ %d: %.2f%%\n", key, percentage)
// 	}

type Engine struct {
	sendCh        []chan Payload
	shard         []map[string]int
	numberOfShard uint32
}

type Payload struct {
	key   string
	value int
}

func (e *Engine) Insert(key string, value int) {
	f := fnv.New32a()
	f.Write([]byte(key))

	shardIndex := int(f.Sum32() % uint32(e.numberOfShard))

	e.sendCh[shardIndex] <- Payload{
		key:   key,
		value: value,
	}
}

func NewEngine(capacity uint32) *Engine {
	if capacity < 1 {
		capacity = 1
	}

	e := Engine{
		sendCh:        make([]chan Payload, capacity),
		shard:         make([]map[string]int, capacity),
		numberOfShard: capacity,
	}

	for index := 0; index < int(capacity); index++ {
		e.sendCh[index] = make(chan Payload, 100) // Buffered channel
		e.shard[index] = make(map[string]int)

		go e.Receive(index)
	}

	return &e
}

func (e *Engine) Receive(index int) {
	for data := range e.sendCh[index] {
		// func()

		e.shard[index][data.key] = data.value
	}
}

func Insert(str map[string]struct{}, data []byte) {
	str[string(data)] = struct{}{}
}
