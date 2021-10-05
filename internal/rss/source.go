package rss

import (
	"encoding/json"
	"sync"

	"github.com/indes/flowerss-bot/internal/config"
	"github.com/indes/flowerss-bot/internal/model"
	"go.uber.org/zap"
)

const (
	SourceTopic   = "source"
	SourceChannel = "sourceChannel"
)

// PubUpdateSources product update source message
func ProduceSourceMessage() {
	if config.RunMode == config.TestMode {
		return
	}
	puber := NewPuber(PuberConfig{
		Host: config.NSQ.NSQd.Host,
		Port: config.NSQ.NSQd.Port,
	})
	sources := model.GetSubscribedNormalSources()
	wg := sync.WaitGroup{}
	for _, source := range sources {
		wg.Add(1)
		go func(s *model.Source) {
			defer wg.Done()
			if !s.NeedUpdate() {
				return
			}
			body, err := json.Marshal(s)
			if err != nil {
				zap.S().Errorw(
					"ProduceSourceMessage, Error to marshal source object to json",
					"source", s,
					"error", err,
				)
			} else {
				zap.S().Debugw(
					"Produce Source Message",
					"source_id", s.ID,
					"source_title", s.Title,
				)
				puber.PubData(SourceTopic, body)
			}
		}(source)
	}
	wg.Wait()
}

func ComsumeSourceMessage(s *model.Source) {
	s.GetNewContents()
}
