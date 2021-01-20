package store

import (
	"fmt"
	"log"
	"time"

	"github.com/18211167516/gocache"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gogf/gf/util/gconv"
)

type MemcachedStore struct {
	MeClient *memcache.Client
	SizeKey  string
}

var MemCachedName = "Memcached"

func NewMemcached(addr string) *MemcachedStore {
	client := memcache.New(addr)

	err := client.Ping()
	if err != nil {
		log.Panic("Cache store Memcached:", err)
	}
	return &MemcachedStore{
		MeClient: client,
		SizeKey:  "size",
	}
}

func RegisterMemcached(addr string) {
	gocache.Register(MemCachedName, NewMemcached(addr))
}

// GC 定期扫删除过期键
func (ms *MemcachedStore) GC() {

}

func (ms *MemcachedStore) GetStoreName() string {
	return MemCachedName
}

func (ms *MemcachedStore) Get(key string) (interface{}, error) {
	item, err := ms.MeClient.Get(key)
	if err != nil {
		return nil, err
	} else {
		return gconv.String(item.Value), err
	}
}

func NewMemItem(key string, value interface{}, d int) *memcache.Item {
	return &memcache.Item{Key: key, Value: gconv.Bytes(value), Expiration: int32(10)}
}

func (ms *MemcachedStore) Incr(key string, i int) error {
	if _, err := ms.MeClient.Increment(key, uint64(i)); err != nil {
		if err == memcache.ErrCacheMiss {
			return ms.Forever(key, i)
		} else {
			return err
		}
	}
	return nil
}

func (ms *MemcachedStore) Decr(key string, i int) error {
	if _, err := ms.MeClient.Decrement(key, uint64(i)); err != nil {
		if err == memcache.ErrCacheMiss {
			return ms.Forever(key, gconv.String(i))
		} else {
			return err
		}
	}
	return nil
}

func (ms *MemcachedStore) Set(key string, value interface{}, d int) error {
	item := NewMemItem(key, value, d)
	if err := ms.MeClient.Set(item); err != nil {
		return err
	} else {
		return ms.Incr(ms.SizeKey, 1)
	}
}

func (ms *MemcachedStore) Delete(key string) error {

	if err := ms.MeClient.Delete(key); err != nil {
		return err
	} else {
		return ms.Decr(ms.SizeKey, 1)
	}

}

func (ms *MemcachedStore) Has(key string) error {
	item, err := ms.MeClient.Get(key)
	if err != nil {
		return err
	}

	if item == nil {
		return fmt.Errorf("Cache store Memcached key:%s not found", key)
	}

	return nil
}

func (ms *MemcachedStore) Forever(key string, value interface{}) error {
	return ms.Set(key, value, 0)
}

func (ms *MemcachedStore) Clear() error {
	return ms.MeClient.DeleteAll()
}

func (ms *MemcachedStore) Size() int {
	value, err := ms.Get(ms.SizeKey)
	if err != nil {
		return 0
	} else {
		return gconv.Int(value) - 1
	}
}

func (ms *MemcachedStore) GetTTl(key string) (time.Duration, error) {
	if err := ms.Has(key); err != nil {
		return 0, err
	} else {
		item, _ := ms.MeClient.Get(key)
		return time.Duration(item.Expiration), nil
	}
}

func (ms *MemcachedStore) IsExpire(key string) (bool, error) {
	if err := ms.Has(key); err != nil {
		return false, err
	} else {
		item, _ := ms.MeClient.Get(key)
		return item.Expiration > 0, nil
	}
}
