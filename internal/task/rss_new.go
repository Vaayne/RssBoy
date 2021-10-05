package task

import (
	"github.com/indes/flowerss-bot/internal/config"
	"github.com/indes/flowerss-bot/internal/rss"
)

func init() {
	task := InitNewRssUpdateTask()
	registerTask(task)
}

// NewRssUpdateTask new NewRssUpdateTask
func InitNewRssUpdateTask() *NewRssUpdateTask {
	return &NewRssUpdateTask{}
}

// NewRssUpdateTask rss source update task
type NewRssUpdateTask struct{}

// Name 任务名称
func (t *NewRssUpdateTask) Name() string {
	return "NewRssUpdateTask"
}

// Stop stop task
func (t *NewRssUpdateTask) Stop() {}

// Start run task
func (t *NewRssUpdateTask) Start() {
	if config.RunMode == config.TestMode {
		return
	}

	// run producer
	go rss.StartSourcePuber()
	go rss.StartContentPuber()

	// run consumer
	go rss.StartSourceComsumer()
	go rss.StartContentComsumer()
}
