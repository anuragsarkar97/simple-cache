package src

import (
	"container/heap"
	"sync"
	"time"
)

type queueData struct {
	key        string
	data       interface{}
	ttl        time.Duration
	expireAt   time.Time
	queueIndex int
}

func newQueueItem(key string, data interface{}, ttl time.Duration) *queueData {
	item := &queueData{
		data: data,
		ttl:  ttl,
		key:  key,
	}
	// since nobody is aware yet of this item, it's safe to touch without lock here
	item.addTimeStamp()
	return item
}

func (q *queueData) expired() bool {
	if q.ttl <= 0 {
		return false
	}
	return q.expireAt.Before(time.Now())
}

func (q *queueData) addTimeStamp() {
	if q.ttl > 0 {
		q.expireAt = time.Now().Add(q.ttl)
	}
}

type priorityQueue struct {
	items []*queueData
}

func newPriorityQueue() *priorityQueue {
	queue := &priorityQueue{}
	heap.Init(queue)
	return queue
}

func (p priorityQueue) Len() int {
	return len(p.items)
}

func (p priorityQueue) Less(i, j int) bool {
	panic("implement me")
}

func (p priorityQueue) Swap(i, j int) {
	panic("implement me")
}

func (p priorityQueue) Push(x interface{}) {
	panic("implement me")
}

func (p priorityQueue) Pop() interface{} {
	panic("implement me")
}

type SimpleCache struct {
	fileName   string
	maxEntry   uint64
	queue      *priorityQueue
	ttl        time.Duration
	data       map[string][]interface{}
	readLock   bool
	lock       *sync.RWMutex
	updateFile chan bool
}

func createNewCache(fileName string, maxEntry uint64) *SimpleCache {
	cache := &SimpleCache{
		fileName: fileName,
		readLock: false,
		maxEntry: maxEntry,
		queue:    newPriorityQueue(),
		ttl:      0,
	}
	go cache.processExpiry()
	return cache
}

func (c *SimpleCache) set(k string, v interface{}) bool {
	// TODO implement set : based on concurrency method
	return true
}

func (c *SimpleCache) get(k string) (interface{}, error, bool) {
	// TODO implement get : based on key available will return result error and presence
	return nil, nil, false
}

func (c *SimpleCache) processExpiry() {
	// TODO check priority queue and keep removing keys that are expired
}

func (c *SimpleCache) close() {
	c.lock.Lock()
	defer c.lock.Unlock()
	// TODO : stop any read write process and update the cache bin file
}

func (c *SimpleCache) setTTL(ttl time.Duration) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.ttl = ttl
}

func (c *SimpleCache) updatePersistentFile() {

}

func (c *SimpleCache) updateReadLock() {
	c.readLock = !c.readLock
}
