package tasks

import (
	"FranceDeveloppe/JEB-backend/tasks/fetcher"
	"time"
)

type backgroundTasks struct {
	Handler func()
	Delay   time.Duration
}

var tasks = []backgroundTasks{
	{
		Handler: SyncUUIDs,
		Delay:   time.Second * 10,
	},
	{
		Handler: fetcher.UpdateData,
		Delay:   time.Second * 30,
	},
}

func RunTasksInBackground() {
	for _, task := range tasks {
		ticker := time.NewTicker(task.Delay)
		go func() {
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					task.Handler()
				}
			}
		}()
	}
}
