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
		Enable:  true,
		Handler: fetcher.UpdateData,
		Delay:   time.Second * 5,
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
			for {
				select {
				case <-ticker.C:
					task.Handler()
				}
			}
		}()
	}
}
