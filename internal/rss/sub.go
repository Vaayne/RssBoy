package rss

import (
	"encoding/json"
	"fmt"

	"github.com/indes/flowerss-bot/internal/config"
	"github.com/indes/flowerss-bot/internal/model"
	"github.com/nsqio/go-nsq"
	"go.uber.org/zap"
)

type SuberConfig struct {
	Host    string
	Port    int32
	Topic   string
	Channel string
}

type Suber struct {
	Cli *nsq.Consumer
}

func NewSuber(c SuberConfig, handler nsq.Handler) (*Suber, error) {
	config := nsq.NewConfig()
	consumer, err := nsq.NewConsumer(
		c.Topic,
		c.Channel,
		config,
	)
	if err != nil {
		zap.S().Errorw(
			"Error create Suber",
			"Topic", c.Topic,
			"Channel", c.Channel,
			"error", err,
		)
		return nil, err
	}
	consumer.AddHandler(handler)
	err = consumer.ConnectToNSQLookupd(fmt.Sprintf("%s:%d", c.Host, c.Port))
	if err != nil {
		zap.S().Errorw(
			"Error ConnectToNSQLookupd",
			"Topic", c.Topic,
			"Channel", c.Channel,
			"Host", c.Host,
			"Port", c.Port,
			"error", err,
		)
		return nil, err
	}
	return &Suber{Cli: consumer}, nil
}

// ContentMessageHandler message handler for subscribe contents
type ContentMessageHandler struct{}

func (h *ContentMessageHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		// Returning nil will automatically send a FIN command to NSQ to mark the message as processed.
		return nil
	}
	var content model.Content
	err := json.Unmarshal(m.Body, &content)
	if err != nil {
		zap.S().Errorw(
			"Error Unmarshal content",
			"data", m.Body,
			"error", err,
		)
	}
	// ComsumeContentMessage
	ComsumeContentMessage(&content)

	// Returning a non-nil error will automatically send a REQ command to NSQ to re-queue the message.
	return nil
}

// SourceMessageHandler message handler for subscribe contents
type SourceMessageHandler struct{}

func (h *SourceMessageHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		// Returning nil will automatically send a FIN command to NSQ to mark the message as processed.
		return nil
	}

	// ComsumeSourceMessage
	var source model.Source
	err := json.Unmarshal(m.Body, &source)
	if err != nil {
		zap.S().Errorw(
			"Error Unmarshal source",
			"data", m.Body,
			"error", err,
		)
	}
	ComsumeSourceMessage(&source)

	// Returning a non-nil error will automatically send a REQ command to NSQ to re-queue the message.
	return nil
}

func StartSourceComsumer() {
	config := SuberConfig{
		Host:    config.NSQ.NSQLookupd.Host,
		Port:    config.NSQ.NSQLookupd.Port,
		Topic:   SourceTopic,
		Channel: SourceChannel,
	}
	suber, err := NewSuber(config, &SourceMessageHandler{})
	if err != nil {
		zap.S().Errorw(
			"Init SourceComsumer Error",
			"config", config,
			"error", err,
		)
	}
	for {
		select {
		case <-suber.Cli.StopChan:
			return
		}
	}
}

func StartContentComsumer() {
	config := SuberConfig{
		Host:    config.NSQ.NSQLookupd.Host,
		Port:    config.NSQ.NSQLookupd.Port,
		Topic:   ContentTopic,
		Channel: ContentChannel,
	}
	suber, err := NewSuber(config, &ContentMessageHandler{})
	if err != nil {
		zap.S().Errorw(
			"Init ContentComsumer Error",
			"config", config,
			"error", err,
		)
	}
	for {
		select {
		case <-suber.Cli.StopChan:
			return
		}
	}
}
