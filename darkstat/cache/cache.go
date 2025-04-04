package cache

import (
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
	first_init  chan error
}

func NewCached[T any](getter func() T, timeToLive time.Duration) *Cached[T] {
	c := &Cached[T]{
		getter:     getter,
		timeToLive: timeToLive,
	}

	c.first_init = async.ToAsync(func() {
		c.get()
	})
	return c
}

func (c *Cached[T]) Get() T {
	if c.first_init != nil {
		<-c.first_init
		c.first_init = nil
	}
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
