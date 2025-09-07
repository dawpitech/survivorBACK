package tasks

import (
	"FranceDeveloppe/JEB-backend/tasks/fetcher"
	"time"
)

type backgroundTasks struct {
	Enable  bool
	Handler func()
	Delay   time.Duration
}

var tasks = []backgroundTasks{
	{
		Enable:  true,
		Handler: SyncUUIDs,
		Delay:   time.Second * 10,
	},
	{
		Enable:  false,
		Handler: fetcher.UpdateData,
		Delay:   time.Minute * 10,
	},
}

func RunTasksInBackground() {
	for _, task := range tasks {
		if !task.Enable {
			continue
		}
		ticker := time.NewTicker(task.Delay)
		go func() {
			defer ticker.Stop()
			task.Handler()
			for {
				select {
				case <-ticker.C:
					task.Handler()
				}
			}
		}()
	}
}
