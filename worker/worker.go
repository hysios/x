package worker

import (
	"errors"
	"reflect"
	"time"

	"github.com/creasty/defaults"
	"github.com/hysios/log"
	"github.com/xxjwxc/gowp/workpool"
)

type JobFunc[T any] func(t T) error

type WorkerOption struct {
	MaxSize      int           `default:"500"`
	Interval     time.Duration `default:"10s"`
	StepInterval time.Duration `default:"10s"`
}

type Worker[T any] struct {
	Option WorkerOption
	iter   Iter[T]
	wp     *workpool.WorkPool
	job    JobFunc[T]
	doJob  func(wp *workpool.WorkPool, u T)
}

type Iter[T any] interface {
	Next() T
	Reset() error
}

func New[K, T any](option WorkerOption, job JobFunc[T]) *Worker[T] {
	defaults.Set(&option)

	return &Worker[T]{
		Option: option,
		// wp:     workpool.New(option.MaxSize),
		job: job,
	}
}

func (worker *Worker[T]) SetIter(iter Iter[T]) {
	worker.iter = iter
}

// SetJob
func (worker *Worker[T]) SetJob(job JobFunc[T]) {
	worker.job = job
}

func (worker *Worker[T]) Iter() Iter[T] {
	return worker.iter
}

func (worker *Worker[T]) isZero(v T) bool {
	vv := reflect.ValueOf(v)

	switch vv.Kind() {
	case reflect.Float32, reflect.Float64, reflect.Int, reflect.Int16,
		reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uint8, reflect.String,
		reflect.Bool:
		return vv.IsZero()
	default:
		if vv.IsZero() || vv.IsNil() {
			return true
		}
	}

	return false
}

func (worker *Worker[T]) Run(loop int) error {
	if worker.iter == nil {
		return errors.New("not set Load Func")
	}

	log.Infof("worker option % #v", worker.Option)
	for i := 0; i != loop; i++ {
		log.Infof("worker loop %d", i)
		worker.wp = workpool.New(worker.Option.MaxSize)
		worker.iter.Reset()
		j := 0
		for v := worker.iter.Next(); !worker.isZero(v); v = worker.iter.Next() {
			func(u T) {
				worker.do(worker.wp, u)
				// worker.wp.Do(func() error {
				// 	t := time.Now()
				// 	if err := worker.job(u); err != nil {
				// 		log.Warnf("worker error %s", err)
				// 		return err
				// 	}

				// 	// 补偿
				// 	s := time.Since(t)
				// 	if s < worker.Option.Interval {
				// 		time.Sleep(worker.Option.Interval - s)
				// 	}

				// 	return nil
				// })

			}(v)
			if (j+1)%worker.Option.MaxSize == 0 {
				log.Infof("worker step wait %s", worker.Option.StepInterval)
				time.Sleep(worker.Option.StepInterval)
			}
			j++
		}
		log.Infof("worker wait %s", worker.Option.Interval)
		time.Sleep(worker.Option.Interval)
		worker.wp.Wait()
	}

	return nil
}

func (worker *Worker[T]) do(wp *workpool.WorkPool, u T) {
	if worker.doJob == nil {
		wp.Do(func() error {
			t := time.Now()
			if err := worker.job(u); err != nil {
				log.Warnf("worker error %s", err)
				return err
			}

			// 补偿
			s := time.Since(t)
			if s < worker.Option.Interval {
				time.Sleep(worker.Option.Interval - s)
			}

			return nil
		})
	} else {
		worker.doJob(wp, u)
	}
}

func (worker *Worker[T]) SetDo(fn func(wp *workpool.WorkPool, u T)) {
	worker.doJob = fn
}
