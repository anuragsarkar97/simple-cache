package main

import (
	"container/heap"
	"sync"
	"time"
)

/*
queue data struct is keep track of individual data point
	1. interface data -> recommended type data []byte
	2. key of the data either string or serialized string
	3. expiry of the key default cache key or custom key
	4. current queue index. higher on the queue denotes closer expiry
*/

type queueData struct {
	key        string
	data       interface{}
	expireAt   int64
	queueIndex int
}

/*
priority queue impl helps in maintaining expiry
initial priority queue length is not fixed but
is configurable through global cache.

Impl will take care of memory size and queue size.
will change dynamically based on the values set in the
cache. 1024 kb to 10,000 keys in current instruction.
*/

type priorityQueue struct {
	items []*queueData
}

/*
Simple cache is a simple cache uses priority queue as
the core dta structure to keep keys.  it uses Least recently
used keys concept to implement key store.
the persistent file is called every 15 minute to update the file. though
this is completely experimental and is computationally expensive
to do so.
It uses read write mutex lock to make sure the read and writes happen without
conflict.
*/

type SimpleCache struct {
	fileName      string
	maxEntry      uint64
	queue         *priorityQueue
	ttl           int64
	data          map[string]*queueData
	lock          *sync.Mutex
	updateFile    chan bool
	expiryChannel bool
}

/*
newQueueItem create a new item to be inserted into the cache.
time stamp is use to set ttl. if tt is -1, it will set to global
cache timing else it will set to the item presented.
*/

func newQueueItem(key string, data interface{}, ttl int64) *queueData {
	item := &queueData{
		data: data,
		key:  key,
	}
	// since nobody is aware yet of this item, it's safe to touch without lock here
	item.addTimeStamp(ttl)
	return item
}

/*
expired checks if the item is expired and removes it form the queue
and the cache server.
*/

func (q *queueData) expired() bool {
	return q.expireAt < time.Now().Unix()
}

/*
addTimeStamp add expiry time to the queue item
will be then used by the processExpiryFunction to remove/add/update
the queue.
*/

func (q *queueData) addTimeStamp(ttl int64) {
	q.expireAt = time.Now().Add(time.Duration(ttl)).Unix()
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
	return p.items[i].expireAt < p.items[j].expireAt
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

func CreateNewCache(fileName string, maxEntry uint64) *SimpleCache {
	cache := &SimpleCache{
		data:     make(map[string]*queueData),
		fileName: fileName,
		maxEntry: maxEntry,
		queue:    newPriorityQueue(),
		ttl:      -1,
		lock:     new(sync.Mutex),
	}
	go cache.concurrentProcessChecks()
	return cache
}

func (c *SimpleCache) Set(k string, v interface{}, ttl int64) bool {
	c.lock.Lock()
	data, present := c.getData(k, ttl)
	if present == true {
		if ttl == -1 {
			data.expireAt = time.Now().Unix() + c.ttl
		} else {
			data.expireAt = time.Now().Unix() + ttl
		}
		c.queue.update(data)
		c.data[k] = data
	} else {
		newData := newQueueItem(k, v, ttl)
		c.queue.push(newData)
		c.data[k] = newData
	}
	c.lock.Unlock()
	c.processExpiry()
	return true
}

func (c *SimpleCache) getData(k string, ttl int64) (*queueData, bool) {
	data, present := c.data[k]
	if !present {
		return nil, false
	}
	if ttl != -1 {
		data.addTimeStamp(ttl)
	} else {
		data.addTimeStamp(c.ttl)
	}
	return data, present
}

func (c *SimpleCache) Get(k string) (interface{}, error, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	data, present := c.getData(k, -1)
	if present {
		return data.data, nil, present
	}
	return nil, nil, false
}

func (c *SimpleCache) processExpiry() {
	c.lock.Lock()
	for c.queue.Len() > 0 && c.queue.items[0].expireAt < time.Now().Unix() {
		delete(c.data, c.queue.items[0].key)
		c.queue.pop()
	}
	c.lock.Unlock()
}

func (c *SimpleCache) concurrentProcessChecks() {
	for c.expiryChannel != true {
		c.processExpiry()
	}
}

func (c *SimpleCache) close() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.updatePersistentFile()
	c.expiryChannel = true
	c = new(SimpleCache)
}

func (c *SimpleCache) setTTL(ttl int64) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.ttl = ttl
}

func (c *SimpleCache) updatePersistentFile() {

}
