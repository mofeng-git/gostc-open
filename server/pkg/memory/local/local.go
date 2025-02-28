package local

import (
	"encoding/json"
	"errors"
	"github.com/patrickmn/go-cache"
	cache_interface "server/pkg/memory"
	"time"
)

type Memory struct {
	db          *cache.Cache
	file        string
	persistence time.Duration
}

func (c *Memory) Get(key string, data cache_interface.AnyType) error {
	bytes, ok := c.db.Get(key)
	if !ok {
		return errors.New("key does not exist")
	}
	temp, ok := bytes.([]byte)
	if !ok {
		return errors.New("value not []byte")
	}
	return data.Unmarshal(temp)
}

func (c *Memory) Set(key string, data cache_interface.AnyType, duration time.Duration) {
	c.db.Set(key, data.Marshal(), duration)
}

func (c *Memory) Del(key string) {
	c.db.Delete(key)
}

func (c *Memory) Sync() {
	_ = c.db.SaveFile(c.file)
}

func NewMemory(file string, duration time.Duration) *Memory {
	c := cache.New(5*time.Minute, 10*time.Minute)
	_ = c.LoadFile(file)
	go func() {
		for {
			// 定期实例化
			time.Sleep(duration)
			_ = c.SaveFile(file)
		}
	}()
	return &Memory{c, file, duration}
}

func (c *Memory) SetStruct(key string, data any, duration time.Duration) {
	marshal, _ := json.Marshal(data)
	c.db.Set(key, marshal, duration)
}

func (c *Memory) GetStruct(key string, data any) error {
	val, ok := c.db.Get(key)
	if !ok {
		return errors.New("key does not exist")
	}
	return json.Unmarshal(val.([]byte), data)
}

func (c *Memory) GetString(key string) (value string) {
	val, ok := c.db.Get(key)
	if !ok {
		return ""
	}
	str, ok := val.(string)
	if ok {
		return str
	}
	return ""
}

func (c *Memory) SetString(key string, value string, duration time.Duration) {
	c.db.Set(key, value, duration)
}
