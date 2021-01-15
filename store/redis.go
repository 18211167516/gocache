package store

import (
	"context"
	"log"
	"time"

	"github.com/18211167516/gocache"
	"github.com/go-redis/redis/v8"
)

var (
	RedisName = "Redis"
	Ctx       context.Context
)

type RedisStore struct {
	reidsClient *redis.Client
}

func NewRedis(addr, password string, db, PoolSize int) *RedisStore {

	ClientRedis := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       db,       // use default DB
		PoolSize: PoolSize, // 连接池大小
		//MinIdleConns: 5,
	})

	Ctx = ClientRedis.Context()

	err := ClientRedis.Ping(Ctx).Err()
	if err != nil {
		log.Fatalln("Cache store redis：", err)
	}

	return &RedisStore{
		reidsClient: ClientRedis,
	}
}

func NewAndRegister(addr, password string, db, PoolSize int) {
	gocache.Register(RedisName, NewRedis(addr, password, db, PoolSize))
}

// GC 定期扫删除过期键
func (ms *RedisStore) GC() {

}

func (ms *RedisStore) GetStoreName() string {
	return RedisName
}

func (ms *RedisStore) Get(key string) (interface{}, error) {
	return ms.reidsClient.Get(Ctx, key).Result()
}

func (ms *RedisStore) Set(key string, value interface{}, d int) error {
	return ms.reidsClient.Set(Ctx, key, value, time.Duration(d)*time.Second).Err()
}

func (ms *RedisStore) Delete(key string) error {
	return ms.reidsClient.Del(Ctx, key).Err()
}

func (ms *RedisStore) Has(key string) error {
	return ms.reidsClient.Exists(Ctx, key).Err()
}

func (ms *RedisStore) Forever(key string, value interface{}) error {
	return ms.reidsClient.Set(Ctx, key, value, 0).Err()
}

func (ms *RedisStore) Clear() error {
	return ms.reidsClient.FlushDB(Ctx).Err()
}

func (ms *RedisStore) Size() int {
	return int(ms.reidsClient.DBSize(Ctx).Val())
}

func (ms *RedisStore) GetTTl(key string) (time.Duration, error) {
	return ms.reidsClient.TTL(Ctx, key).Result()
}

func (ms *RedisStore) IsExpire(key string) (bool, error) {
	res, err := ms.reidsClient.Exists(Ctx, key).Result()
	if err != nil {
		return false, err
	} else {
		if res > 0 {
			return false, err
		} else {
			return true, err
		}
	}
}
