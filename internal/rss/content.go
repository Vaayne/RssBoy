package rss

import (
	"encoding/json"
	"sync"

	"github.com/indes/flowerss-bot/internal/bot"
	"github.com/indes/flowerss-bot/internal/config"
	"github.com/indes/flowerss-bot/internal/model"
	"go.uber.org/zap"
)

const (
	ContentTopic   = "content"
	ContentChannel = "contentChannel"
)

func ProduceContentMessage() {
	puber := NewPuber(PuberConfig{
		Host: config.NSQ.NSQd.Host,
		Port: config.NSQ.NSQd.Port,
	})

	contents := model.QueryUnBroadcastedContent()
	wg := sync.WaitGroup{}
	for _, content := range contents {
		wg.Add(1)
		go func(c model.Content) {
			defer wg.Done()
			body, err := json.Marshal(c)
			if err != nil {
				zap.S().Errorw(
					"Error to marshal source object to json",
					"content", c,
					"error", err,
				)
			} else {
				zap.S().Debugw(
					"Publish Content Message",
					"title", c.Title,
					"hash_id", c.HashID,
				)
				puber.PubData(ContentTopic, body)
			}
		}(content)
	}
	wg.Wait()
}

func ComsumeContentMessage(content *model.Content) {
	source, err := model.GetSourceById(content.SourceID)
	if err != nil {
		zap.S().Errorw(
			"Error GetSourceById",
			"ID", content.SourceID,
			"error", err,
		)
		return
	}
	subs := model.GetSubscriberBySource(source)
	bot.BroadcastContent(source, subs, content)
	model.SetIsBroadTrue(content)
}
