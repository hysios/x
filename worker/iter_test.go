package worker

import (
	"errors"
	"log"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type User struct {
	ID        uint   `faker:"-"`
	Username  string `faker:"username"`
	FirstName string `faker:"first_name"`
	Age       int    `faker:"boundary_start=15,boundary_end=80"`
}

var db *gorm.DB

func testDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatalf("open sqlite memory failed")
	}

	// db.Exec("DELETE FROM `users`")
	db.AutoMigrate(&User{})

	for i := 0; i < 100; i++ {
		var u User
		err := faker.FakeData(&u)
		if err != nil {
			log.Printf("faker data error %s", err)
		}
		db.Create(&u)
	}
	return db
}

func loadUsers(id uint, size int) (last uint, users []*User, err error) {
	users = make([]*User, 0)
	log.Printf("load from %d", id)
	if err := db.Limit(size).Find(&users, "id > ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrEmptySlice) {
			return 0, nil, ErrEndIter
		}
	}
	if len(users) == 0 {
		return 0, nil, ErrEndIter
	}
	lastU := users[len(users)-1]
	if len(users) < size {
		return lastU.ID, users, ErrEndIter
	}

	return lastU.ID, users, nil
}

func TestNewDBIter(t *testing.T) {
	db = testDB()
	const size = 10

	var c = 0
	iter := NewDBIter(func(id uint) (uint, []*User, error) {
		var users = make([]*User, 0)
		t.Logf("from id %d", id)
		if err := db.Limit(size).Find(&users, "id > ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrEmptySlice) {
				return 0, nil, ErrEndIter
			}
		}
		c++
		if len(users) == 0 {
			return 0, nil, ErrEndIter
		}
		last := users[len(users)-1]

		if len(users) < size {
			return last.ID, users, ErrEndIter
		}

		return last.ID, users, nil
	}, 0)

	i := 0
	for u := iter.Next(); u != nil; u = iter.Next() {
		t.Logf("user %v", u)
		i++
	}

	assert.Equal(t, 100, i)
	assert.Equal(t, 11, c)
}
