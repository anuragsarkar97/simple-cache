package main

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	c := CreateNewCache("", 1000)
	c.setTTL(40)
	for i := 0; i < 100; i++ {
		k := RandStringRunes(5)
		v := RandStringRunes(8)
		c.Set(k, v, -1)
	}
	fmt.Println("set")
	time.Sleep(1 * time.Second)
	fmt.Println(c.queue.Len())
	c.close()

}

var letterRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_#@!$%")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
