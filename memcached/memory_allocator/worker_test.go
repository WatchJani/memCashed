package memory_allocator

import (
	"bytes"
	"fmt"
	"log"
	"testing"

	"github.com/WatchJani/memCashed/memcached/parser"
)

func Test(t *testing.T) {
	allocator := New(5 * 1024 * 1024)

	slabsSize := []int{
		64, 128, 256,
		512, 1024, 2048,
		4096, 8192, 16384,
		32768, 65536, 131072,
		262144, 524288, 1048576,
	}

	slabAllocator := make([]Slab, len(slabsSize))
	for i := range slabAllocator {
		slabAllocator[i] = NewSlab(slabsSize[i], 0, allocator)
	}

	slabManager := NewSlabManager(slabAllocator, 1)

	setData := []string{
		"dJNhxBA", "dJNhxBA",
		"ghmifXY", "ghmifXY",
		"aVtYnis", "aVtYnis",
		"iooJPrV", "iooJPrV",
		"CfZBRRJ", "CfZBRRJ",
		"ftgNWRj", "ftgNWRj",
		"MGcXBVO", "MGcXBVO",
		"SAIExZG", "SAIExZG",
		"tWrByDf", "tWrByDf",
		"NeZvjTf", "NeZvjTf",
	}
	var writer *bytes.Buffer = &bytes.Buffer{}
	for index := 0; index < len(setData); index += 2 {
		payload, err := parser.Encode('S', []byte(setData[index]), []byte(setData[index+1]), -1)
		if err != nil {
			log.Println(err)
		}

		slabManager.chooseOperation(NewTransfer(payload[4:], 0, writer))
	}

	slabManager.store.Range(func(key, value interface{}) bool {
		fmt.Println(string(value.(Key).field))
		return true
	})
}
