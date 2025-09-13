package cache

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/darklab8/fl-darkstat/darkstat/settings/logus"
	"github.com/darklab8/go-utils/typelog"
	"github.com/darklab8/go-utils/utils/async"
	"github.com/darklab8/go-utils/utils/ptr"
)

type Cached[T any] struct {
	value        *T
	getter       func() T
	timeToLive   time.Duration
	timeCreated  time.Time
	first_init   sync.WaitGroup
	ExtraLogging bool
}

var Lock sync.Mutex

type CacheOption[T any] func(c *Cached[T])

func NewCached[T any](getter func() T, timeToLive time.Duration, opts ...CacheOption[T]) *Cached[T] {
	c := &Cached[T]{
		getter:     getter,
		timeToLive: timeToLive,
	}
	c.timeCreated = time.Now()

	for _, opt := range opts {
		opt(c)
	}

	go func() {
		for {
			c.first_init.Add(1)
			async.ToAsync(func() {
				Lock.Lock()
				defer Lock.Unlock()
				c.get()
				c.first_init.Done()
				logus.Log.Debug("updated cache with time to live", typelog.Any("ttl_period", c.timeToLive.Seconds()))
			})
			time.Sleep(timeToLive - time.Second*20)
		}
	}()
	return c
}

func (c *Cached[T]) Get(ctx context.Context) T {
	c.first_init.Wait()
	return c.get()
}

func (c *Cached[T]) get() T {
	expiry_date := c.timeCreated.Add(c.timeToLive)
	if c.ExtraLogging {
		fmt.Println("CACHED CACHED CACHED: expire date=", expiry_date, " now data=", time.Now())
	}
	if c.value == nil {
		if c.ExtraLogging {
			fmt.Println("CACHED CACHED CACHED: is nil, updated")
		}
		c.value = ptr.Ptr(c.getter())
		c.timeCreated = time.Now()
		logus.Log.Debug("CACHED nil cache calced ")
	} else if time.Now().After(expiry_date) {
		if c.ExtraLogging {
			fmt.Println("CACHED CACHED CACHED: after expired data succeeded, updating")
		}
		c.value = ptr.Ptr(c.getter())
		c.timeCreated = time.Now()
		logus.Log.Debug("CACHED is expired and recalced")
	}
	logus.Log.Debug("CACHED is returned", typelog.Float64("ttl_left", expiry_date.Sub(time.Now()).Seconds()))
	return *c.value
}
