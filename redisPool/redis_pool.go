package redisPool

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/jericho-yu/aid/dict"
	redis "github.com/redis/go-redis/v9"
)

type (
	RedisPool struct {
		connections *dict.AnyDict[string, *redisConn]
	}

	redisConn struct {
		prefix string
		conn   *redis.Client
	}
)

var (
	redisPoolIns  *RedisPool
	redisPoolOnce sync.Once
	RedisPoolApp  RedisPool
)

// Once 单例化：redis 链接
func (*RedisPool) Once(redisSetting *RedisSetting) *RedisPool {
	redisPoolOnce.Do(func() {
		redisPoolIns = &RedisPool{}
		redisPoolIns.connections = dict.MakeAnyDict[string, *redisConn]()

		if len(redisSetting.Pool) > 0 {
			for _, pool := range redisSetting.Pool {
				redisPoolIns.connections.Set(pool.Key, &redisConn{
					prefix: fmt.Sprintf("%s:%s", redisSetting.Prefix, pool.Prefix),
					conn: redis.NewClient(&redis.Options{
						Addr:     fmt.Sprintf("%s:%d", redisSetting.Host, redisSetting.Port),
						Password: redisSetting.Password,
						DB:       pool.DbNum,
					}),
				})
			}
		}
	})

	return redisPoolIns
}

// GetClient 获取链接和链接前缀
func (*RedisPool) GetClient(key string) (string, *redis.Client) {
	if client, exist := redisPoolIns.connections.Get(key); exist {
		return client.prefix, client.conn
	}
	return "", nil
}

// Get 获取值
func (*RedisPool) Get(clientName, key string) (string, error) {
	var (
		err         error
		prefix, ret string
		client      *redis.Client
	)

	prefix, client = redisPoolIns.GetClient(clientName)
	if client == nil {
		return "", fmt.Errorf("没有找到redis链接：%s", clientName)
	}

	ret, err = client.Get(context.Background(), fmt.Sprintf("%s:%s", prefix, key)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		} else {
			return "", err
		}
	}
	return ret, nil
}

// Set 设置值
func (*RedisPool) Set(clientName, key string, val any, exp time.Duration) (string, error) {
	var (
		prefix string
		client *redis.Client
	)

	prefix, client = redisPoolIns.GetClient(clientName)
	if client == nil {
		return "", fmt.Errorf("没有找到redis链接：%s", clientName)
	}

	return client.Set(context.Background(), fmt.Sprintf("%s:%s", prefix, key), val, exp).Result()
}

// Close 关闭链接
func (my *RedisPool) Close(key string) error {
	if client, exist := redisPoolIns.connections.Get(key); exist {
		return client.conn.Close()
	}
	return nil
}

// Clean 清理链接
func (*RedisPool) Clean() {
	for key, val := range redisPoolIns.connections.All() {
		val.conn.Close()
		redisPoolIns.connections.RemoveByKey(key)
	}
}
