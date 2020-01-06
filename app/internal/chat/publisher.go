package chat

import (
	"encoding/json"
	"github.com/go-redis/redis/v7"
)

type Publisher struct {
	client  *redis.Client
	channel string
}

func (p *Publisher) Publish(message Message) {
	b, err := json.Marshal(message)
	err = p.client.Publish(p.channel, string(b)).Err()
	if err != nil {
		panic(err)
	}
}

func NewPublisher(client *redis.Client, channel string) Publisher {
	return Publisher{
		client:  client,
		channel: channel,
	}
}
