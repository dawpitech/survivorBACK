package tasks

import (
	"FranceDeveloppe/JEB-backend/tasks/fetcher"
	"time"
)

func RunTasksInBackground() {
	ticker := time.NewTicker(time.Second * 10)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				SyncUUIDs()
				fetcher.UpdateData()
			}
		}
	}()
}
