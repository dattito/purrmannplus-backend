package scheduler

import (
	"time"

	"github.com/dattito/purrmannplus-backend/utils/logging"
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

// Start the scheduler async
func StartAsync() {
	logging.Info("Starting scheduler async")
	S.StartAsync()
}

// Start the scheduler sync
func StartBlocking() {
	logging.Info("Starting scheduler blocking")
	S.StartBlocking()
}
