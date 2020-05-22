package main

import (
	"container/heap"
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

var (
	wg sync.WaitGroup
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

/*
Queue Data struct is keep track of individual Data point
	1. interface Data -> recommended type Data []byte
	2. Key of the Data either string or serialized string
	3. expiry of the Key default cache Key or custom Key
	4. current Queue index. higher on the Queue denotes closer expiry
*/

type queueData struct {
	Key        string      `json:"key"`
	Data       interface{} `json:"data"`
	ExpireAt   int64       `json:"expire_at"`
	QueueIndex int         `json:"index"`
}

/*
priority Queue impl helps in maintaining expiry
initial priority Queue length is not fixed but
is configurable through global cache.

Impl will take care of memory size and Queue size.
will change dynamically based on the values set in the
cache. 1024 kb to 10,000 keys in current instruction.
*/

type priorityQueue struct {
	Items []*queueData `json:"queue_data"`
}

/*
Simple cache is a simple cache uses priority Queue as
the core dta structure to keep keys.  it uses Least recently
used keys concept to implement Key store.
the persistent file is called every 15 minute to update the file. though
this is completely experimental and is computationally expensive
to do so.
It uses read write mutex Lock to make sure the read and writes happen without
conflict.
*/

type SimpleCache struct {
	FileName      string                `json:"file_name"`
	MaxEntry      uint64                `json:"max_entry"`
	Queue         *priorityQueue        `json:"queue"`
	TTL           int64                 `json:"cache_global_ttl"`
	Data          map[string]*queueData `json:"data"`
	Lock          *sync.Mutex           `json:"lock"`
	ExpiryChannel chan bool             `json:"-"`
	SaveFile      bool                  `json:"-"`
}

/*
newQueueItem create a new item to be inserted into the cache.
time stamp is use to set TTL. if tt is -1, it will set to global
cache timing else it will set to the item presented.
*/

func newQueueItem(key string, data interface{}, ttl int64) *queueData {
	item := &queueData{
		Data: data,
		Key:  key,
	}
	// since nobody is aware yet of this item, it's safe to touch without Lock here
	item.addTimeStamp(ttl)
	return item
}

/*
expired checks if the item is expired and removes it form the Queue
and the cache server.
*/

func (q *queueData) expired() bool {
	return q.ExpireAt < time.Now().Unix()
}

/*
addTimeStamp add expiry time to the Queue item
will be then used by the processExpiryFunction to remove/add/update
the Queue.
*/

func (q *queueData) addTimeStamp(ttl int64) {
	q.ExpireAt = time.Now().Add(time.Duration(ttl)).Unix()
}

func (p *priorityQueue) update(data *queueData) {
	heap.Fix(p, data.QueueIndex)
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
	heap.Remove(p, queueData.QueueIndex)
}

func (p *priorityQueue) Len() int {
	return len(p.Items)
}

func (p *priorityQueue) Less(i, j int) bool {
	return p.Items[i].ExpireAt < p.Items[j].ExpireAt
}

func (p *priorityQueue) Swap(i, j int) {
	p.Items[i], p.Items[j] = p.Items[j], p.Items[i]
	p.Items[i].QueueIndex = i
	p.Items[j].QueueIndex = j
}

func (p *priorityQueue) Push(x interface{}) {
	item := x.(*queueData)
	item.QueueIndex = len(p.Items)
	p.Items = append(p.Items, item)
}

func (p *priorityQueue) Pop() interface{} {
	old := p.Items
	n := len(old)
	item := old[n-1]
	item.QueueIndex = -1
	p.Items = old[0 : n-1]
	return item
}

func newPriorityQueue() *priorityQueue {
	queue := &priorityQueue{}
	heap.Init(queue)
	return queue
}

func CreateNewCache(fileName string, maxEntry uint64, save bool) *SimpleCache {
	cache := &SimpleCache{
		Data:          make(map[string]*queueData),
		FileName:      fileName,
		MaxEntry:      maxEntry,
		Queue:         newPriorityQueue(),
		TTL:           -1,
		Lock:          new(sync.Mutex),
		ExpiryChannel: make(chan bool),
		SaveFile:      save,
	}
	wg.Add(1)
	go cache.concurrentProcessChecks()
	return cache
}

func (c *SimpleCache) Set(k string, v interface{}, ttl int64) bool {
	c.Lock.Lock()
	data, present := c.getData(k, ttl)
	if present == true {
		if ttl == -1 {
			data.ExpireAt = time.Now().Unix() + c.TTL
		} else {
			data.ExpireAt = time.Now().Unix() + ttl
		}
		c.Queue.update(data)
		c.Data[k] = data
	} else {
		var ttx int64
		if ttl == -1 {
			ttx = time.Now().Unix() + c.TTL
		} else {
			ttx = time.Now().Unix() + ttl
		}
		newData := newQueueItem(k, v, ttx)
		c.Queue.push(newData)
		c.Data[k] = newData
	}
	c.Lock.Unlock()
	return true
}

func (c *SimpleCache) getData(k string, ttl int64) (*queueData, bool) {
	data, present := c.Data[k]
	if !present {
		return nil, false
	}
	if ttl != -1 {
		data.addTimeStamp(ttl)
	} else {
		data.addTimeStamp(c.TTL)
	}
	return data, present
}

func (c *SimpleCache) Get(k string) (interface{}, error, bool) {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	data, present := c.getData(k, -1)
	if present {
		return data.Data, nil, present
	}
	return nil, nil, false
}

func (c *SimpleCache) processExpiry() {
	c.Lock.Lock()
	for c.Queue.Len() > 0 && c.Queue.Items[0].ExpireAt < time.Now().Unix() {
		delete(c.Data, c.Queue.Items[0].Key)
		c.Queue.pop()
	}
	c.Lock.Unlock()
}

func (c *SimpleCache) concurrentProcessChecks() {
	defer wg.Done()
	for {
		select {
		case <-c.ExpiryChannel:
			c.ExpiryChannel <- true
			return
		default:
			break
		}
		c.processExpiry()
	}
}

func (c *SimpleCache) close() {
	c.ExpiryChannel <- true
	<-c.ExpiryChannel
	wg.Wait()
	close(c.ExpiryChannel)
	if c.SaveFile {
		c.updatePersistentFile()
	}
}

func (c *SimpleCache) setTTL(ttl int64) {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	c.TTL = ttl
}

func (c *SimpleCache) updatePersistentFile() {
	x, _ := json.Marshal(c)
	err := writeGob(c.FileName, string(x))
	if err != nil {
		log.Printf("error occured while because : %s\n", err.Error())
	}
}

func writeGob(filePath string, data string) error {
	file, err := os.Create(filePath)
	file.WriteString(data)
	file.Close()
	return err
}
