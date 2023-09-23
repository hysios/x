package cache

import (
	"time"
)

type FetchOption struct {
	RefreshTTL int64
}

type UpdateOption struct {
	TTL int64
}

type LoadOpt func(*FetchOption)

type UpdateOpt func(*UpdateOption)

// Alive
func (opt *FetchOption) Alive() time.Duration {
	return time.Duration(opt.RefreshTTL) * time.Second
}

// TTLDuration
func (opt *UpdateOption) TTLDuration() time.Duration {
	return time.Duration(opt.TTL) * time.Second
}

// WithTTL
func ResetTTL(ttl int64) LoadOpt {
	return func(opt *FetchOption) {
		opt.RefreshTTL = ttl
	}
}

// WithTTL
func WithTTL(ttl int64) UpdateOpt {
	return func(opt *UpdateOption) {
		opt.TTL = ttl
	}
}
