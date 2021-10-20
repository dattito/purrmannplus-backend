package scheduler

import (
	"time"

	"github.com/go-co-op/gocron"
)

var S *gocron.Scheduler

// Initizalize the scheduler object
func Init() {
	S = gocron.NewScheduler(time.Local)
}

// Add a job to the scheduler object
func AddJob(cron string, exec func()) {
	S.Cron(cron).Do(exec)
}

// Start the scheduler
func StartScheduler() {
	S.StartAsync()
}
