package redis

import (
	"context"
	"fmt"

	"github.com/hysios/x/cache"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type CacheOption struct {
	Encoder   cache.Encoder
	Decoder   cache.Decoder
	TTL       int64
	Namespace string
	Log       *zap.Logger
}

type CacheOpt func(*CacheOption)

func New[Key, Value any](redisCli *redis.Client, opts ...CacheOpt) cache.Cache[Key, Value] {
	var opt = &CacheOption{
		Encoder:   cache.DefaultEncoder,
		Decoder:   cache.DefaultDecoder,
		Namespace: cache.Namespace,
		Log:       zap.NewNop(),
	}

	for _, o := range opts {
		o(opt)
	}

	return &redisCache[Key, Value]{
		cli:       redisCli,
		enc:       opt.Encoder,
		dec:       opt.Decoder,
		namespace: opt.Namespace,
		ttl:       opt.TTL,
		log:       opt.Log,
	}
}

type redisCache[Key, Value any] struct {
	cli *redis.Client

	prefix    string
	enc       cache.Encoder
	dec       cache.Decoder
	ttl       int64
	namespace string
	log       *zap.Logger
}

// key
func (r *redisCache[Key, Value]) key(key Key) string {
	return fmt.Sprintf("%s:%s%v", r.namespace, r.prefix, key)
}

func (r *redisCache[Key, Value]) Load(key Key, opts ...cache.LoadOpt) (val Value, ok bool) {
	var (
		ctx = context.Background()
		opt = &cache.FetchOption{}
	)

	for _, o := range opts {
		o(opt)
	}

	var (
		v   string
		err error
	)
	if opt.Alive() > 0 {
		v, err = r.cli.GetEx(ctx, r.key(key), opt.Alive()).Result()
	} else {
		v, err = r.cli.Get(ctx, r.key(key)).Result()
	}

	if err != nil {
		ok = false
		return
	}

	if err = r.dec.Unmarshal([]byte(v), &val); err != nil {
		ok = false
		return
	}

	ok = true
	return
}

func (r *redisCache[Key, Value]) Update(key Key, val Value, opts ...cache.UpdateOpt) {
	var (
		ctx = context.Background()
		opt = &cache.UpdateOption{
			TTL: r.ttl,
		}
	)

	for _, o := range opts {
		o(opt)
	}

	data, err := r.enc.Marshal(val)
	if err != nil {
		return
	}

	r.log.Debug("redis set", zap.String("key", r.key(key)), zap.String("value", string(data)), zap.Int64("ttl", opt.TTL))
	_, err = r.cli.Set(ctx, r.key(key), data, opt.TTLDuration()).Result()
	if err != nil {
		r.log.Warn("redis set error", zap.Error(err))
		return
	}
}

func (r *redisCache[Key, Value]) Clear(key Key) {
	var ctx = context.Background()
	_ = r.cli.Del(ctx, r.key(key)).Err()
}
