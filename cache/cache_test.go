package cache

import (
	"bytes"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"testing"
)

const N = 1000

var getTests = []struct {
	key   string
	value []byte
}{
	{"123456", []byte("test1")},
	{"12", []byte("test2")},
}

func TestGetSet(t *testing.T) {
	c := New()
	for _, tt := range getTests {
		c.Set(tt.key, tt.value)
		val, err := c.Get(tt.key)
		if err != nil {
			t.Fatalf(err.Error())
		}
		if !bytes.Equal(tt.value, val) {
			t.Fatalf("cache hit = %v; want %v", val, tt.value)
		}
	}
}

func TestGetSetConcurrent(t *testing.T) {
	runtime.GOMAXPROCS(2)

	c := New()
	ints := rand.Perm(N)

	var wg sync.WaitGroup
	wg.Add(len(ints))
	for i := 0; i < len(ints); i++ {
		go func(i int) {
			c.Set(fmt.Sprintf("%d", i), []byte(fmt.Sprintf("%d", i)))
			wg.Done()
		}(i)
	}

	wg.Wait()
	for _, i := range ints {
		if _, err := c.Get(fmt.Sprintf("%d", i)); err != nil {
			t.Errorf(err.Error())
		}
	}
}
