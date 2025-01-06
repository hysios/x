package worker

import "errors"

var (
	ErrEndIter  = errors.New("End Iter")
	ErrStopIter = errors.New("Stop Iter")
	ErrSkipIter = errors.New("Skip Iter")
)
