package redis

import (
	"log"
	"time"

	"github.com/garyburd/redigo/redis"
)

// Cache is Redis cache adapter.
type RedisCache struct {
	p        *redis.Pool // redis connection pool
	conninfo string
	dbNum    int
	password string
}

// actually do the redis cmds
func (rc *RedisCache) do(commandName string, args ...interface{}) (reply interface{}, err error) {
	c := rc.p.Get()
	defer c.Close()

	return c.Do(commandName, args...)
}

// Get cache from redis.
func (rc *RedisCache) Get(key string) interface{} {
	v, err := rc.do("GET", key)
	if err == nil {
		return v
	}
	log.Printf("redis get cache error:%v\n", err)
	return nil
}

// Put put cache to redis.
func (rc *RedisCache) Put(key string, val interface{}, timeout time.Duration) error {
	var err error
	if _, err = rc.do("SETEX", key, int64(timeout/time.Second), val); err != nil {
		return err
	}

	return err
}

// Delete delete cache in redis.
func (rc *RedisCache) Delete(key string) error {
	var err error
	if _, err = rc.do("DEL", key); err != nil {
		return err
	}
	// _, err = rc.do("HDEL", rc.key, key)
	return err
}

// IsExist check cache's existence in redis.
func (rc *RedisCache) IsExist(key string) bool {
	v, err := redis.Bool(rc.do("EXISTS", key))
	if err != nil {
		return false
	}
	// if !v {
	// 	if _, err = rc.do("HDEL", rc.key, key); err != nil {
	// 		return false
	// 	}
	// }
	return v
}

// connect to redis.
func (rc *RedisCache) connectInit() {
	dialFunc := func() (c redis.Conn, err error) {
		c, err = redis.Dial("tcp", rc.conninfo)
		if err != nil {
			return nil, err
		}

		if rc.password != "" {
			if _, err := c.Do("AUTH", rc.password); err != nil {
				c.Close()
				return nil, err
			}
		}

		_, selecterr := c.Do("SELECT", rc.dbNum)
		if selecterr != nil {
			c.Close()
			return nil, selecterr
		}
		return
	}
	// initialize a new pool
	rc.p = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 180 * time.Second,
		Dial:        dialFunc,
	}
}

func RedisNewCache(constr string, password string) *RedisCache {
	rc := RedisCache{}

	rc.conninfo = constr
	rc.dbNum = 0
	rc.password = password

	rc.connectInit()

	return &rc
}
