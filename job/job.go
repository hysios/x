package job

import (
	"errors"
	"os"
	"os/signal"
	"time"

	"github.com/gomodule/redigo/redis"
	v9redis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/gocraft/work"
)

type Worker struct {
	Pool *work.WorkerPool
	cfg  *Config

	start bool
}

type (
	Job          = work.Job
	ScheduledJob = work.ScheduledJob
)

var Namespace = "mx_worker"

type Config struct {
	Concurrency int
	Namespace   string
	Context     interface{}
	RedisPool   *redis.Pool
}

type Context struct {
	DB  *gorm.DB
	Rdb *v9redis.Client
}

var Default = NewWorker(&Config{
	Concurrency: 10,
	Namespace:   Namespace,
	Context:     Context{},
})

// NewWorker creates a new worker
func NewWorker(cfg *Config) *Worker {
	var redisPool = &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", ":6379")
		},
	}

	if cfg.Context == nil {
		cfg.Context = Context{}
	}

	if cfg.RedisPool != nil {
		redisPool = cfg.RedisPool
	} else {
		cfg.RedisPool = redisPool
	}

	pool := work.NewWorkerPool(cfg.Context, uint(cfg.Concurrency), cfg.Namespace, redisPool)

	return &Worker{
		Pool: pool,
		cfg:  cfg,
	}
}

// Start starts the worker
func (w *Worker) Start() error {
	w.start = true
	w.Pool.Start()
	return nil
}

// Stop stops the worker
func (w *Worker) Stop() error {
	w.start = false
	w.Pool.Stop()
	return nil
}

// IsStarted
func (w *Worker) IsStarted() bool {
	return w.start
}

// Enqueue
func (w *Worker) Enqueue(jobName string, jobData map[string]interface{}) (*Job, error) {
	var enqueuer = work.NewEnqueuer(w.cfg.Namespace, w.cfg.RedisPool)

	return enqueuer.Enqueue(jobName, jobData)
}

// EnqueueIn
func (w *Worker) EnqueueIn(jobName string, delay time.Duration, jobData map[string]interface{}) (*ScheduledJob, error) {
	var enqueuer = work.NewEnqueuer(w.cfg.Namespace, w.cfg.RedisPool)
	return enqueuer.EnqueueIn(jobName, int64(delay.Seconds()), jobData)
}

// EnqueueAt
func (w *Worker) EnqueueAt(jobName string, ts time.Time, jobData map[string]interface{}) (*ScheduledJob, error) {
	var (
		now = time.Now()
	)
	if ts.Before(now) {
		return nil, errors.New("job time is in the past")
	}

	return w.EnqueueIn(jobName, time.Since(now), jobData)
}

// Cron
func (w *Worker) Cron(spec string, jobName string) {
	w.Pool.PeriodicallyEnqueue(spec, jobName)
}

// RegisterJob registers a job
func (w *Worker) RegisterJob(jobName string, fn interface{}) {
	if jobName == "" {
		panic("job name is required")
	}
	w.Pool.Job(jobName, fn)
}

// Use
func (w *Worker) Use(middlewares ...interface{}) {
	for _, m := range middlewares {
		w.Pool.Middleware(m)
	}
}

// Start
func Start() error {
	return Default.Start()
}

// SetNamespace
func SetNamespace(namespace string) {
	Namespace = namespace
}

// SetDB
func (ctx *Context) SetDB(db *gorm.DB) {
	ctx.DB = db
}

// SetRdb
func (ctx *Context) SetRdb(rdb *v9redis.Client) {
	ctx.Rdb = rdb
}

func init() {
	go func() {
		HandleExit(Default)
	}()
}

func HandleExit(w *Worker) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan
	if w.IsStarted() {
		w.Stop()
	}
}
