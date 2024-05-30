package redis

import "go.uber.org/zap"

// WithLogger
func WithLogger(log *zap.Logger) CacheOpt {
	return func(opt *CacheOption) {
		opt.Log = log
	}
}

// WithNamespace
func WithNamespace(ns string) CacheOpt {
	return func(opt *CacheOption) {
		opt.Namespace = ns
	}
}

func WithTTL(ttl int64) CacheOpt {
	return func(opt *CacheOption) {
		opt.TTL = ttl
	}
}

// WithKeyGen
func WithKeyGen(fn func(key interface{}) string) CacheOpt {
	return func(opt *CacheOption) {
		opt.KeyGen = fn
	}
}
