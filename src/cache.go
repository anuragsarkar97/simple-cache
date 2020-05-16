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

func (p *priorityQueue) update(data *queueData) {
	heap.Fix(p, data.queueIndex)
}

func (p *priorityQueue) push(data *queueData) {
	heap.Push(p, data)
}

func (p *priorityQueue) pop() *queueData {
	if p.Len() == 0 {
		return nil
	}
	return heap.Pop(p).(*queueData)
}

func (p *priorityQueue) remove(queueData *queueData) {
	heap.Remove(p, queueData.queueIndex)
}

func (p *priorityQueue) Len() int {
	return len(p.items)
}

func (p *priorityQueue) Less(i, j int) bool {
	if p.items[i].expireAt.IsZero() {
		return false
	}
	if p.items[j].expireAt.IsZero() {
		return true
	}
	return p.items[i].expireAt.Before(p.items[j].expireAt)
}

func (p *priorityQueue) Swap(i, j int) {
	p.items[i], p.items[j] = p.items[j], p.items[i]
	p.items[i].queueIndex = i
	p.items[j].queueIndex = j
}

func (p *priorityQueue) Push(x interface{}) {
	item := x.(*queueData)
	item.queueIndex = len(p.items)
	p.items = append(p.items, item)
}

func (p *priorityQueue) Pop() interface{} {
	old := p.items
	n := len(old)
	item := old[n-1]
	item.queueIndex = -1
	p.items = old[0 : n-1]
	return item
}

func newPriorityQueue() *priorityQueue {
	queue := &priorityQueue{}
	heap.Init(queue)
	return queue
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
