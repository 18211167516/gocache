package store

import (
	"container/list"
	"fmt"
	"sync"
	"time"

	"github.com/18211167516/gocache"
)

type MemoryStore struct {
	list       map[string]*list.Element
	expireList *list.List
	mu         sync.RWMutex
}

type Memory struct {
	key     string
	value   interface{}
	expirat time.Time
}

var name = "Memory"

func init() {
	gocache.Register(name, NewStore())
}

func NewStore() *MemoryStore {
	return &MemoryStore{
		list:       make(map[string]*list.Element),
		expireList: list.New(),
	}
}

func (m *Memory) Get() interface{} {
	return m.value
}

func (m *Memory) Set(key string, value interface{}, d int) {
	m.key = key
	m.value = value
	m.expirat = time.Now().Add(time.Duration(d) * time.Second)
}

func (m *Memory) Forever(key string, value interface{}) {
	m.key = key
	m.value = value
}

func (m *Memory) TTL(key string) time.Time {
	return m.expirat
}

// GC 定期扫删除过期键
func (ms *MemoryStore) GC() {
	ticker := time.NewTicker(2 * time.Second)
	for {
		fmt.Println("GC Start")
		select {
		case <-ticker.C:
			// 触发定时器
			elememt := ms.expireList.Back()
			if elememt == nil {
				break
			}

			key := elememt.Value.(*Memory).key
			if isbool, _ := ms.IsExpire(key); isbool {
				ms.expireList.Remove(elememt)
				delete(ms.list, key)
				fmt.Println("缓存删除成功:", key)
			}

		}

	}

}

func (ms *MemoryStore) GetStoreName() string {
	return name
}

func (ms *MemoryStore) Get(key string) (interface{}, error) {
	defer ms.mu.RUnlock()
	ms.mu.RLock()
	isExpre, err := ms.IsExpire(key)

	if err != nil {
		return nil, err
	} else {
		if isExpre {
			ms.expireList.Remove(ms.list[key])
			delete(ms.list, key)
			return nil, fmt.Errorf("Cache :%s not found", key)
		} else {
			return ms.list[key].Value.(*Memory).Get(), nil
		}
	}

}

func (ms *MemoryStore) Set(key string, value interface{}, time int) error {
	if value == nil {
		return fmt.Errorf("Cache : %s set value not is nil", "Memory")
	}
	ms.mu.Lock()
	defer ms.mu.Unlock()
	m := &Memory{}
	m.Set(key, value, time)
	element := ms.expireList.PushBack(m)
	ms.list[key] = element
	return nil
}

func (ms *MemoryStore) Delete(key string) error {
	defer ms.mu.Unlock()
	ms.mu.Lock()
	if _, ok := ms.list[key]; !ok {
		return fmt.Errorf("Cache :%s value not found", key)
	} else {
		ms.expireList.Remove(ms.list[key])
		delete(ms.list, key)
		return nil
	}
}

func (ms *MemoryStore) Has(key string) error {
	defer ms.mu.RUnlock()
	ms.mu.RLock()
	if _, ok := ms.list[key]; ok {
		return nil
	} else {
		return fmt.Errorf("Cache :%s value not found", key)
	}
}

func (ms *MemoryStore) Forever(key string, value interface{}) error {
	defer ms.mu.Unlock()
	ms.mu.Lock()
	if value == nil {
		return fmt.Errorf("Cache : %s set value not is nil", "Memory")
	}
	m := &Memory{}
	m.Forever(key, value)
	//无过期时间的只存在list map
	list := list.New()
	element := list.PushBack(m)
	ms.list[key] = element
	return nil
}

func (ms *MemoryStore) Clear() error {
	defer ms.mu.Unlock()
	ms.mu.Lock()
	ms.list = make(map[string]*list.Element)
	ms.expireList = list.New()
	return nil
}

func (ms *MemoryStore) Size() int {
	defer ms.mu.RUnlock()
	ms.mu.RLock()
	return len(ms.list)
}

func (ms *MemoryStore) GetTTl(key string) (time.Time, error) {
	defer ms.mu.RUnlock()
	ms.mu.RLock()
	if err := ms.Has(key); err != nil {
		return time.Time{}, err
	} else {
		return ms.list[key].Value.(*Memory).TTL(key), nil
	}
}

func (ms *MemoryStore) IsExpire(key string) (bool, error) {
	expire, err := ms.GetTTl(key)
	if err != nil {
		return false, err
	}

	return !expire.IsZero() && time.Now().After(expire), nil
}
