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

// WithPrefix
