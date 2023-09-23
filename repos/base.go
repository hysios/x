package repos

import (
	"reflect"

	"github.com/hysios/x/providers"

	"gorm.io/gorm"
)

type Repos[Record any] interface {
	isRepos()
}

type Base[Record any, Key any] interface {
	// Create a new record
	Create(t *Record) error
	// Get a record by id
	Get(id Key) (*Record, error)
	// Find all records
	FindAll() ([]*Record, error)
	// Update a record
	Update(t *Record) error
	// Delete a record
	Delete(id Key) error

	isRepos()
}

type initer interface {
	init()
}

// Init a repo
func Init[Model any, Key any, R Repos[Model]](db *gorm.DB) R {
	var r R
	t := reflect.TypeOf(new(R)).Elem()

	baseCtor, ok := impls.Lookup(t)
	if !ok {
		panic("not have Repos[Record] implements")
	}

	var ctor func(db *gorm.DB) R
	if ctor, ok = baseCtor.(func(db *gorm.DB) R); !ok {
		panic("not have Repos[Record] implements")
	}

	r = ctor(db)

	if i, ok := any(r).(initer); ok {
		i.init()
	}

	extends, ok := registers.Lookup(t)
	if !ok {
		panic("not have Repos[Record] implements")
	}

	for _, ext := range extends {
		// t := reflect.TypeOf(r).Elem()
		r = ext.(Extender[Model, R])(r, db)
	}

	return r
	// return &base[Record, Key]{db: db.(gorm.DB)}
}

type Injector[Record any, R Repos[Record]] func(r Repos[Record])

type Extender[Record any, R Repos[Record]] func(base R, db *gorm.DB) R

func Extend[Record any, R Repos[Record]](ext Extender[Record, R]) {
	t := reflect.TypeOf(new(R)).Elem()
	registers.Register(t, ext)
}

var registers providers.SliceProvider[reflect.Type, any]

// Impl
func Impl[Record any, R Repos[Record]](ctor func(db *gorm.DB) R) {
	t := reflect.TypeOf(new(R)).Elem()
	impls.Register(t, ctor)
}

var impls providers.Provider[reflect.Type, any]
