package datastructure

import (
	"sync"
)

type CSlice struct {
	sync.RWMutex
	size    int
	cap     int
	entries []interface{}
}

type Iterator struct {
	Ind   int
	Entry interface{}
}

func CSliceObj(icap int) *CSlice {
	slice := &CSlice{
		size:    0,
		cap:     icap,
		entries: make([]interface{}, 0),
	}
	return slice
}

func (slice *CSlice) Append(entry interface{}) {
	slice.Lock()
	defer slice.Unlock()
	if slice.size == slice.cap {
		slice.entries = slice.entries[1:]
		slice.size--
	}
	slice.entries = append(slice.entries, entry)
	slice.size++
}

func (slice *CSlice) List() <-chan Iterator {
	channel := make(chan Iterator)

	go func(slice *CSlice) {
		slice.Lock()
		defer slice.Unlock()
		for i, e := range slice.entries {
			iter := Iterator{i, e}
			channel <- iter
		}
		close(channel)
	}(slice)

	return channel
}
