package rediswrapper

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type RedisCli struct {
	Cli  *redis.Client
	Ctx  context.Context
	Conf *Config
}

func New(ctx context.Context, conf Config) *RedisCli {
	cli := redis.NewClient(&redis.Options{DB: conf.DB})
	return &RedisCli{
		Cli: cli,
		Ctx: ctx,
	}
}

func (r *RedisCli) Get(key string) (string, error) {
	val, err := r.Cli.Get(r.Ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", errors.Wrap(err, "Redis cache missed for key")
	}
	if err != nil {
		return "", errors.Wrap(err, "failed to get from Redis")
	}
	var answers string
	if err := json.Unmarshal([]byte(val), &answers); err != nil {
		return "", errors.Wrap(err, "failed to unmarshal value from Redis")
	}
	return answers, nil
}
func (r *RedisCli) Set(key, val string, ttl int) error {
	valJson, err := json.Marshal(val)
	if err != nil {
		return errors.Wrap(err, "failed to marshal value")
	}
	dur := time.Second * time.Duration(ttl)
	valStr := string(valJson)
	if _, err := r.Cli.Set(r.Ctx, key, valStr, dur).Result(); err != nil {
		log.Println("failed to set redis", key, val)
	}
	return nil
}
