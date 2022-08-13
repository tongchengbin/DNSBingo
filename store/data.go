package store

import (
	"errors"
	"github.com/patrickmn/go-cache"
	"sync"
	"time"
)

type Storage struct {
	cache *cache.Cache
}

// CorrelationData is the data for a correlation-id.
type CorrelationData struct {
	// data contains data for a correlation-id in AES encrypted json format.
	Data []string `json:"data"`
	// dataMutex is a mutex for the data slice.
	dataMutex *sync.Mutex
}

func (c *Storage) RegisterKey(key string) error {
	_, found := c.cache.Get(key)
	if found {
		return errors.New("key provided already exists")
	}
	data := &CorrelationData{
		Data:      make([]string, 0),
		dataMutex: &sync.Mutex{},
	}
	c.cache.Set(key, data, 10*time.Minute)
	return nil
}

func (c *Storage) SetItem(key string, val string) error {
	item, found := c.cache.Get(key)
	if !found {
		return errors.New("could not get correlation-id from cache")
	}
	value, ok := item.(*CorrelationData)
	if !ok {
		return errors.New("invalid correlation-id cache value found")
	}
	value.dataMutex.Lock()
	value.Data = append(value.Data, val)
	value.dataMutex.Unlock()
	return nil
}

func (c *Storage) GetItem(key string) (*CorrelationData, error) {
	item, ok := c.cache.Get(key)
	if !ok {
		return nil, errors.New("cache item not found")
	}
	value, ok := item.(*CorrelationData)
	if !ok {
		return nil, errors.New("cache item not found")
	}
	return value, nil
}

// New creates a new storage instance for interactsh data.
func New() *Storage {
	c := cache.New(5*time.Minute, 10*time.Minute)
	return &Storage{cache: c}
}

func (c *Storage) SetData(key string, val interface{}) {
	c.cache.Set(key, val, 10*time.Minute)
}

func (c *Storage) GetData(key string) (interface{}, bool) {
	return c.cache.Get(key)
}

var Store = New()
