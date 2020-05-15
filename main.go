package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	c := Cache{}
	c.INIT()
	for i := 0; i < 100; i++ {
		k := RandStringRunes(5)
		v := RandStringRunes(8)
		c.Set(k, v)
	}
	for k := range c.data {
		fmt.Println(k, c.data[k])
	}
}

type Cache struct {
	concurrent bool
	data       map[string]string
	readLock   bool
	lock       *sync.Mutex
}

func (c *Cache) INIT() *Cache {
	c.data = make(map[string]string)
	c.lock = new(sync.Mutex)
	c.concurrent = false
	return c
}

func (c *Cache) Get(k string) (string, error) {
	if c.readLock != true {
		return c.data[k], nil
	}
	if c.concurrent == true {
		return "", errors.New("write lock applied")
	} else {
		for c.readLock == true {
			time.Sleep(10)
			if c.readLock == false {
				break
			}
		}
		return c.data[k], nil
	}
}

func (c *Cache) Set(k string, v string) {
	if c.concurrent == true {
		go setData(k, v, &c.data, c.lock, &c.readLock)
	} else {
		setData(k, v, &c.data, c.lock, &c.readLock)
	}
}

func setData(k string, v string, d *map[string]string, lock *sync.Mutex, readLock *bool) {
	*readLock = true
	lock.Lock()
	m := *d
	m[k] = v
	lock.Unlock()
	*readLock = false
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_#@!$%")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
