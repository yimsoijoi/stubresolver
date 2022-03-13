package cacher

import "github.com/yimsoijoi/stubresolver/rediswrapper"

type Cacher struct {
	Redis *rediswrapper.RedisCli
}

func New(r *rediswrapper.RedisCli) *Cacher {
	return &Cacher{
		Redis: r,
	}
}
