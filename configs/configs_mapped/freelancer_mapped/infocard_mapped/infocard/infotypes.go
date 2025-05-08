package infocard

import "sync"

type Infocard struct {
	content string
	Lines   []string
}

func NewInfocard(content string) *Infocard {
	return &Infocard{content: content}
}

type Infoname string

type RecordKind string

const (
	TYPE_NAME    RecordKind = "NAME"
	TYPE_INFOCAD RecordKind = "INFOCARD"
)

type Config struct {
	infonames map[int]Infoname
	infocards map[int]*Infocard
	mutex     *sync.RWMutex
	unsafe    *UnsafeConfig
}

func NewConfig() *Config {
	c := &Config{
		infocards: make(map[int]*Infocard),
		infonames: make(map[int]Infoname),
		mutex:     new(sync.RWMutex),
	}
	c.unsafe = &UnsafeConfig{
		Infonames: c.infonames,
		Infocards: c.infocards,
		Mutex:     c.mutex,
	}
	return c
}

func (c *Config) GetInfocard(id int) *Infocard {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.infocards[id]
}
func (c *Config) GetInfoname(id int) Infoname {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.infonames[id]
}
func (c *Config) GetInfocard2(id int) (*Infocard, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	value, ok := c.infocards[id]
	return value, ok
}
func (c *Config) GetInfoname2(id int) (Infoname, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	value, ok := c.infonames[id]
	return value, ok
}
func (c *Config) PutInfocard(id int, value *Infocard) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.infocards[id] = value
}
func (c *Config) PutInfoname(id int, value Infoname) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.infonames[id] = value
}

func (c *Config) GetDicts(callback func(config *UnsafeConfig)) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	callback(c.unsafe)
}

func (c *Config) GetUnsafe() *UnsafeConfig {
	return c.unsafe
}

type UnsafeConfig struct {
	Infonames map[int]Infoname
	Infocards map[int]*Infocard
	Mutex     *sync.RWMutex
}
