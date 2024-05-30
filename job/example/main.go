package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/hysios/x/job"
)

func main() {
	var worker = job.NewWorker(&job.Config{
		Namespace:   "test",
		Concurrency: 10,
		Context:     job.Context{},
	})

	worker.RegisterJob("echo", func(job *job.Job) error {
		log.Printf("echo %s", job.ArgString("message"))
		return nil
	})

	worker.Enqueue("echo", map[string]interface{}{"message": "hello"})

	go func() {
		var tick = time.NewTicker(5 * time.Second)
		for range tick.C {
			worker.EnqueueIn("echo", 10*time.Second, map[string]interface{}{"message": "hello world"})
		}
	}()

	worker.Cron("0 0 * * * *", "echo")

	worker.Start()
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

	// Stop the pool
	worker.Stop()
}
