package action

import (
	"errors"
)

var (
	ErrStop  = errors.New("stop")
	ErrRetry = errors.New("retry")
)

// type errRetry struct {
// 	err error
// }

// func (e *errRetry) Error() string {
// 	return e.err.Error()
// }

// func (e *errRetry) Unwrap() error {
// 	return e.err
// }

// func (e *errRetry) Is(target error) bool {
// 	_, ok := target.(*errRetry)
// 	return ok
// }

// func ErrRetry(err error) error {
// 	return &errRetry{err: err}
// }
