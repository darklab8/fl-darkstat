package cache

import (
	"sync"
	"time"

	"github.com/darklab8/fl-darkstat/darkstat/settings/logus"
	"github.com/darklab8/go-typelog/typelog"
	"github.com/darklab8/go-utils/utils/async"
	"github.com/darklab8/go-utils/utils/ptr"
)

type Cached[T any] struct {
	value       *T
	getter      func() T
	timeToLive  time.Duration
	timeCreated time.Time
	first_init  sync.WaitGroup
}

var Lock sync.Mutex

func NewCached[T any](getter func() T, timeToLive time.Duration) *Cached[T] {
	c := &Cached[T]{
		getter:     getter,
		timeToLive: timeToLive,
	}

	go func() {
		c.first_init.Add(1)
		async.ToAsync(func() {
			Lock.Lock()
			defer Lock.Unlock()
			c.get()
			c.first_init.Done()
		})
		time.Sleep(timeToLive - time.Second*20)
	}()

	return c
}

func (c *Cached[T]) Get() T {
	c.first_init.Wait()
	return c.get()
}

func (c *Cached[T]) get() T {
	expiry_date := c.timeCreated.Add(c.timeToLive)
	if c.value == nil {
		c.value = ptr.Ptr(c.getter())
		c.timeCreated = time.Now()
		logus.Log.Debug(" nil cache calced ")
	} else if time.Now().After(expiry_date) {
		c.value = ptr.Ptr(c.getter())
		c.timeCreated = time.Now()
		logus.Log.Debug("cache is expired and recalced")
	}
	logus.Log.Debug("cache is returned", typelog.Float64("ttl_left", expiry_date.Sub(time.Now()).Seconds()))
	return *c.value
}
