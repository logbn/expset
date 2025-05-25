package expset

import (
	"sync"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/elliotchance/orderedmap/v3"
)

type item[T comparable] struct {
	exp int64
	ttl time.Duration
	val T
}

// Set implements a generic expiring set.
type Set[T comparable] struct {
	clock clock.Clock
	done  chan bool
	exp   *orderedmap.OrderedMap[int64, T]
	items *orderedmap.OrderedMap[T, *item[T]]
	mutex sync.RWMutex
	wg    sync.WaitGroup
}

// New returns a new [Cache].
func New[T comparable]() *Set[T] {
	return &Set[T]{
		clock: clock.New(),
		done:  make(chan bool),
		exp:   orderedmap.NewOrderedMap[int64, T](),
		items: orderedmap.NewOrderedMap[T, *item[T]](),
	}
}

// Start begins the eviction process.
func (c *Set[T]) Start() {
	c.wg.Add(1)
	tick := c.clock.Ticker(time.Second)
	go func() {
		defer tick.Stop()
		for {
			select {
			case t := <-tick.C:
				c.mutex.Lock()
				for exp, v := range c.exp.AllFromFront() {
					if exp > t.UnixNano() {
						break
					}
					c.evict(v)
				}
				c.mutex.Unlock()
			case <-c.done:
				c.wg.Done()
				return
			}
		}
	}()
}

// Stop stops the eviction process.
func (c *Set[T]) Stop() {
	close(c.done)
	c.wg.Wait()
}

// Add adds an item to the set with a given time to live.
func (c *Set[T]) Add(val T, ttl time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	i, ok := c.items.Get(val)
	if ok {
		c.exp.Delete(i.exp)
	} else {
		i = &item[T]{val: val}
		c.items.Set(val, i)
	}
	i.ttl = ttl
	i.exp = c.clock.Now().Add(ttl).UnixNano()
	for c.exp.Has(i.exp) {
		i.exp++
	}
	c.exp.Set(i.exp, val)
}

// Refresh refreshes the expiration based on the original TTL
func (c *Set[T]) Refresh(val T) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	i, ok := c.items.Get(val)
	if !ok {
		return false
	}
	exp := i.exp
	i.exp = c.clock.Now().Add(i.ttl).UnixNano()
	for c.exp.Has(i.exp) {
		i.exp++
	}
	c.exp.Delete(exp)
	c.exp.Set(i.exp, i.val)
	return true
}

// Has indicates whether the set contains a value.
func (c *Set[T]) Has(val T) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.items.Has(val)
}

// Clear clears the set.
func (c *Set[T]) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.exp = orderedmap.NewOrderedMap[int64, T]()
	c.items = orderedmap.NewOrderedMap[T, *item[T]]()
}

// Len returns the number of items in the set.
func (c *Set[T]) Len() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.items.Len()
}

func (c *Set[T]) evict(val T) bool {
	i, _ := c.items.Get(val)
	c.items.Delete(val)
	c.exp.Delete(i.exp)
	return true
}
