package main

import "github.com/indes/flowerss-bot/internal/rss"

func main() {

	// run consumer
	go rss.StartSourceComsumer()
	go rss.StartContentComsumer()
}
