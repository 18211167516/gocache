package gocache

import (
	"fmt"
	"log"
	"strings"
	"time"
)

var Stores = make(map[string]StoreInterface)

type Cache struct {
	name  string
	store StoreInterface
}

type StoreInterface interface {
	GetStoreName() string
	// 获取缓存
	Get(key string) (interface{}, error)
	// 设置缓存带过期时间
	Set(key string, value interface{}, time int) error
	// 设置永久缓存无过期时间
	Forever(key string, value interface{}) error
	// 删除key
	Delete(key string) error
	// 判断key是否存在
	Has(key string) error
	// 全部清空
	Clear() error
	// 获取缓存key的数量
	Size() int
	// 获取expire
	GetTTl(key string) (time.Time, error)
	// 随机删除已过期key
	GC()
}

// 获取一个实例
func New(name string) (*Cache, error) {
	name = strings.ToLower(name)
	if store, ok := Stores[name]; ok {
		go store.GC()
		return &Cache{
			name:  name,
			store: store,
		}, nil
	} else {
		return nil, fmt.Errorf("Cache:unknown %s store please import", name)
	}

}

// 注册store
func Register(name string, store StoreInterface) {
	if store == nil {
		log.Panic("Cache: Register store is nil")
	}

	name = strings.ToLower(name)

	if _, ok := Stores[name]; ok {
		log.Panic("Cache: Register store is exist")
	}

	Stores[name] = store
}

// 获取store名称
func (c *Cache) GetStoreName() string {
	return c.store.GetStoreName()
}

// 获取键值
func (c *Cache) Get(key string) (interface{}, error) {
	return c.store.Get(key)
}

// 设置键值以及过期时间
func (c *Cache) Set(key string, value interface{}, time int) error {
	return c.store.Set(key, value, time)
}

// 永久设置键值不过期
func (c *Cache) Forever(key string, value interface{}) error {
	return c.store.Forever(key, value)
}

// 删除键值
func (c *Cache) Delete(key string) error {
	return c.store.Delete(key)
}

// 判断键是否存在
func (c *Cache) Has(key string) error {
	return c.store.Has(key)
}

// 清除所有键值
func (c *Cache) Clear() error {
	return c.store.Clear()
}

func (c *Cache) Size() int {
	return c.store.Size()
}

func (c *Cache) GetTTl(key string) (time.Time, error) {
	return c.store.GetTTl(key)
}
