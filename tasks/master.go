package tasks

import (
	"FranceDeveloppe/JEB-backend/tasks/fetcher"
	"time"
)

type backgroundTasks struct {
	Enable     bool
	Repetitive bool
	Handler    func()
	Delay      time.Duration
}

var tasks = []backgroundTasks{
	{
		Enable:     true,
		Repetitive: true,
		Handler:    SyncUUIDs,
		Delay:      time.Second * 10,
	},
	{
		Enable:     true,
		Repetitive: false,
		Handler:    fetcher.UpdateData,
		Delay:      time.Minute * 10,
	},
	{
		Enable:     true,
		Repetitive: true,
		Handler:    UpdateUsersWithoutPP,
		Delay:      time.Second * 30,
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
			if !task.Repetitive {
				return
			}
			for {
				select {
				case <-ticker.C:
					task.Handler()
				}
			}
		}()
	}
}
