package main

import (
	"container/heap"
	"sync"
	"time"
)

type Data struct {
	mu             *sync.RWMutex
	sets           map[string]*SetEntry
	expirySetsHeap *SetsHeap
}

type SetsHeap []*SetEntry

type SetEntry struct {
	expiry    time.Duration
	createdAt time.Time
	index     int
	key       string
	data      []byte
}

type SetArgs struct {
	expiry time.Duration
}

func (sh SetsHeap) Len() int {
	return len(sh)
}

func (sh SetsHeap) Less(i int, j int) bool {
	return sh[i].createdAt.Add(sh[i].expiry).Before(sh[j].createdAt.Add(sh[j].expiry))
}

func (sh SetsHeap) Swap(i, j int) {
	sh[i], sh[j] = sh[j], sh[i]
	sh[i].index = i
	sh[j].index = j
}

func (sh *SetsHeap) Push(x interface{}) {
	entry := x.(*SetEntry)
	entry.index = len(*sh)
	*sh = append(*sh, entry)
}

func (sh *SetsHeap) Pop() interface{} {
	old := *sh
	n := len(old)
	x := old[n-1]
	x.index = -1
	*sh = old[:n-1]
	return x
}

func newData(reapInterval time.Duration) *Data {
	sh := &SetsHeap{}

	heap.Init(sh)

	data := &Data{
		mu:             &sync.RWMutex{},
		sets:           make(map[string]*SetEntry),
		expirySetsHeap: sh,
	}

	if reapInterval != 0 {
		go data.reapLoop(reapInterval)
	}

	return data
}

func (d *Data) getSetData(key string) ([]byte, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	se, ok := d.sets[key]
	if !ok {
		return []byte{}, ok
	}
	if se.isExpired() {
		return []byte{}, false
	}
	return se.data, ok
}

func (se *SetEntry) isExpired() bool {
	if se.expiry == 0 {
		return false
	}
	return time.Now().UTC().After(se.createdAt.Add(se.expiry))
}

func (d *Data) setSetData(key string, value []byte, setArgs SetArgs) {
	entry := &SetEntry{
		key:       key,
		data:      value,
		expiry:    setArgs.expiry,
		createdAt: time.Now(),
	}

	d.mu.RLock()
	if oldEntry, exists := d.sets[key]; exists {
		heap.Remove(d.expirySetsHeap, oldEntry.index)
	}
	d.mu.RUnlock()

	d.mu.Lock()
	defer d.mu.Unlock()

	d.sets[key] = entry

	if entry.expiry != 0 {
		heap.Push(d.expirySetsHeap, entry)
	}
}

func (d *Data) reapLoop(reapInterval time.Duration) {
	ticker := time.NewTicker(reapInterval)
	for range ticker.C {
		d.reap()
	}
}

func (d *Data) reap() {
	d.mu.Lock()
	defer d.mu.Unlock()
	for d.expirySetsHeap.Len() > 0 {
		if top := (*d.expirySetsHeap)[0]; top.isExpired() {
			heap.Pop(d.expirySetsHeap)
			delete(d.sets, top.key)
		} else {
			break
		}
	}
}
