package cacher

import (
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type Key struct {
	Dom  string
	Type string
	TTL  int
}

// d domainname + "."(root server) | t type | ttl expire time |
func NewKey(d, t string, ttl int) Key {
	if ttl < 0 {
		ttl = 60
	}
	if d[len(d)-1] != '.' {
		d = d + "."
	}
	return Key{
		Dom:  d,
		Type: t,
		TTL:  ttl,
	}
}

func (k Key) RedisKey() string {
	return k.Dom
}

func (c *Cacher) Get(k Key) ([]string, error) {
	val, err := c.Redis.Get(k.RedisKey())
	if errors.Is(err, redis.Nil) {
		log.Println("Redis cache miss", k.RedisKey())
	}
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (c *Cacher) Set(k Key, val string, ttl Key) error {
	if err := c.Redis.Set(k.RedisKey(), val, ttl.TTL); err != nil {
		return errors.Wrap(err, "failed to set Redis")
	}
	return nil
}
