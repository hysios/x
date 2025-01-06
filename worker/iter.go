package worker

import (
	"errors"

	"github.com/fatih/structs"
	"gorm.io/gorm"
)

type DBIter[K, T any] struct {
	load   LoadFunc[K, T]
	idx    K
	inital K
	remind []T
}

type LoadFunc[K, T any] func(id K) (last K, rets []T, err error)

func NewDBIter[K, T any](load LoadFunc[K, T], inital K) *DBIter[K, T] {
	return &DBIter[K, T]{load: load, idx: inital, inital: inital}
}

func (dbiter *DBIter[K, T]) Next() T {
	var z T
	if len(dbiter.remind) == 0 {
		if last, loads, err := dbiter.load(dbiter.idx); err != nil {
			if len(loads) == 0 {
				return z
			} else {
				dbiter.remind = loads
				dbiter.idx = last
			}
		} else {
			dbiter.remind = loads
			dbiter.idx = last
		}
	}

	next := dbiter.remind[0]
	dbiter.remind = dbiter.remind[1:]
	return next
}

func (dbiter *DBIter[K, T]) Reset() error {
	dbiter.remind = nil
	dbiter.idx = dbiter.inital
	return nil
}

func LoadModel[K, T any](db *gorm.DB, inital K, step int) LoadFunc[K, T] {
	return func(id K) (last K, rets []T, err error) {
		var z K
		rets = make([]T, 0)
		if err := db.Limit(step).Find(&rets, "id > ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrEmptySlice) {
				return z, nil, ErrEndIter
			}
		}

		if len(rets) == 0 {
			return z, nil, ErrEndIter
		}

		l := rets[len(rets)-1]
		if len(rets) < step {
			return objId[K](l), rets, ErrEndIter
		}
		return objId[K](l), rets, err
	}
}

func objId[T any](val any) T {
	s := structs.New(val)
	return s.Field("ID").Value().(T)
}
