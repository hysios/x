package worker

type RedisIter[T any] struct {
	load   LoadFunc[int64, T]
	idx    int64
	inital int64
	remind []T
}

func NewRedisIter[T any](load LoadFunc[int64, T], inital int64) *RedisIter[T] {
	return &RedisIter[T]{load: load, idx: inital, inital: inital}
}

func (riter *RedisIter[T]) Next() T {
	var z T
	if len(riter.remind) == 0 {
		last, loads, err := riter.load(riter.idx)
		if err != nil {
			return z
		} else if last == 0 {
			loads = append(loads, z)
		}

		riter.remind = loads
		riter.idx = last
	}

	next := riter.remind[0]
	riter.remind = riter.remind[1:]
	return next
}

func (riter *RedisIter[T]) Reset() error {
	riter.idx = 0
	riter.inital = 0
	riter.remind = nil

	return nil
}
