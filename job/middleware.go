package job

import (
	"github.com/gocraft/work"
	v9redis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func MiddlewareDB(db *gorm.DB) func(job *work.Job, next work.NextMiddlewareFunc) error {
	return func(job *work.Job, next work.NextMiddlewareFunc) error {
		job.Args["db"] = db
		return next()
	}
}

func MiddlewareRedis(redis *v9redis.Client) func(job *work.Job, next work.NextMiddlewareFunc) error {
	return func(job *work.Job, next work.NextMiddlewareFunc) error {
		job.Args["redis"] = redis
		return next()
	}
}
