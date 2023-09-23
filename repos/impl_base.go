package repos

import (
	"gorm.io/gorm"
)

type fmter interface {
	NoFmt() string
	SetNo(string)
}

type BaseImpl[Record any, Key any] struct {
	DB *gorm.DB

	plugins Plugins
}

// init
func (b *BaseImpl[Record, Key]) init() {
	b.plugins = basePlugins
}

func (b *BaseImpl[Record, Key]) Create(t *Record) error {

	if err := b.plugins.Create(t); err != nil {
		return err
	}

	return b.DB.Create(t).Error
}

func (b *BaseImpl[Record, Key]) Get(id Key) (*Record, error) {
	panic("nonimplement")
}

func (b *BaseImpl[Record, Key]) FindAll() ([]*Record, error) {
	panic("nonimplement")
}

func (b *BaseImpl[Record, Key]) Update(t *Record) error {
	panic("nonimplement")
}

func (b *BaseImpl[Record, Key]) Delete(id Key) error {
	panic("nonimplement")
}

// isRepos()
func (b *BaseImpl[Record, Key]) isRepos() {}

// generateNoFor
func generateNoFor(t interface{}, format string) (string, error) {
	return "", nil
}

var basePlugins Plugins

func init() {
	// basePlugins.AddField()
}
