package repos

import (
	"fmt"

	"github.com/hysios/x/events"
	"gorm.io/gorm"
)

type User struct {
	Id   uint
	No   string
	Name string
}

// Pinyin
func (*User) Pinyin() string {
	return "YongHu"
}

// SetNo set no
func (u *User) SetNo(no string) {
	u.No = no
}

type UserRepos interface {
	Base[User, uint]
}

type userRepos struct {
	Base[User, uint]

	db *gorm.DB
}

type uesrEventRepos struct {
	Base[User, uint]

	bus *events.Bus
}

// Create
func (u *uesrEventRepos) Create(t *User) error {
	// u.bus.Send("user.create", t)
	return u.Base.Create(t)
}

func init() {
	Impl[User](func(db *gorm.DB) Base[User, uint] {
		return &BaseImpl[User, uint]{DB: db}
	})

	Extend[User](func(base Base[User, uint], db *gorm.DB) Base[User, uint] {
		return &userRepos{Base: base, db: db}
	})

	Extend[User](func(base Base[User, uint], db *gorm.DB) Base[User, uint] {
		return &uesrEventRepos{Base: base}
	})
}

func ExampleRepo() {
	db, _ := openDB()
	var userRepo = Init[User, uint, Base[User, uint]](db)
	var user = &User{Name: "张三"}
	userRepo.Create(user)
	fmt.Println(user.No)
	// Output: ZS
	//

}
