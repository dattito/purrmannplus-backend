package scheduler

import (
	"time"

	"github.com/go-co-op/gocron"
)

var S *gocron.Scheduler

func Init() {
	S = gocron.NewScheduler(time.Local)
}

func AddJob(cron string, exec func()) {
	S.Cron(cron).Do(exec)
}

func StartScheduler() {
	S.StartAsync()
}
