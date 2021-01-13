package store

import (
	"fmt"
	"sync"
	"time"
)

type MemoryStore struct {
	name string
	list map[string]*Memory
	lock sync.Mutex
}

type Memory struct {
	value   interface{}
	expirat time.Duration
}

func (m *Memory) Get() interface{} {
	return m.value
}

func (m *Memory) Set(value interface{}, time time.Duration) {
	m.value = value
	m.expirat = time
}

func (ms *MemoryStore) GetStoreName() string {
	return ms.name
}

func (ms *MemoryStore) Get(key string) (interface{}, error) {
	if store, ok := ms.list[key]; ok {
		return store.Get(), nil
	} else {
		return nil, fmt.Errorf("Cache :%s not found", key)
	}
}

func (ms *MemoryStore) Pull(key string) (interface{}, error) {
	if store, ok := ms.list[key]; ok {
		value := store.Get()
		if err := ms.Delete(key); err != nil {
			return value, err
		} else {
			return value, nil
		}
	} else {
		return nil, fmt.Errorf("Cache :%s not found", key)
	}
}

func (ms *MemoryStore) Set(key string, value interface{}, time time.Duration) error {
	if value == nil {
		return fmt.Errorf("Cache : %s set value not is nil", "Memory")
	}
	m := new(Memory)
	m.Set(value, time)
	ms.list[key] = m
	return nil
}

func (ms *MemoryStore) Delete(key string) error {
	if err := ms.Has(key); err != nil {
		return err
	} else {
		delete(ms.list, key)
		return nil
	}
}

func (ms *MemoryStore) Has(key string) error {
	if _, ok := ms.list[key]; ok {
		return nil
	} else {
		return fmt.Errorf("Cache :%s value not found", key)
	}
}

func (ms *MemoryStore) Forever(key string, value interface{}) error {
	if value == nil {
		return fmt.Errorf("Cache : %s set value not is nil", "Memory")
	}
	m := new(Memory)
	m.Set(value, 0)
	ms.list[key] = m
	return nil
}

func (ms *MemoryStore) Clear() error {
	ms.list = make(map[string]*Memory)
	return nil
}
