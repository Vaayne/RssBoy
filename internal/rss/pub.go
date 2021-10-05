package rss

import (
	"fmt"
	"time"

	"github.com/indes/flowerss-bot/internal/config"
	"github.com/nsqio/go-nsq"
	"go.uber.org/zap"
)

type PuberConfig struct {
	Host string
	Port int32
}

type Puber struct {
	Cli *nsq.Producer
}

func NewPuber(c PuberConfig) *Puber {
	config := nsq.NewConfig()
	client, err := nsq.NewProducer(
		fmt.Sprintf("%s:%d", c.Host, c.Port),
		config,
	)
	if err != nil {
		zap.S().Errorw(
			"Error create Puber",
			"Host", c.Host,
			"port", c.Port,
			"error", err,
		)
	}
	return &Puber{
		Cli: client,
	}
}

func (p *Puber) PubData(topic string, body []byte) error {
	err := p.Cli.Publish(topic, body)
	if err != nil {
		zap.S().Errorw(
			"Error pub data",
			"Topic", topic,
			"Data", body,
			"error", err,
		)
		return err
	}
	return nil
}

func StartSourcePuber() {
	for {
		ProduceSourceMessage()
		time.Sleep(time.Duration(config.UpdateInterval) * time.Minute)
	}
}

func StartContentPuber() {
	for {
		ProduceContentMessage()
		time.Sleep(time.Duration(1) * time.Minute)
	}
}
