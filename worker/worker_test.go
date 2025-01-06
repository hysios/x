package worker

import (
	"testing"
)

func TestWorker(t *testing.T) {
	db = testDB()

	worker := New[*User](WorkerOption{}, func(u *User) error {
		t.Logf("user %v", u)
		return nil
	})

	worker.SetIter(NewDBIter[uint, *User](func(id uint) (last uint, rets []*User, err error) {
		return loadUsers(id, 10)
	}, 0))
	worker.Run(1)
}
