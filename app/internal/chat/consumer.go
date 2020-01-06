package chat

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type Consumer struct {
	client        *redis.Client
	connectionMap *ConnectionMap
	pubsub        *redis.PubSub
	channel       string
}

func (c *Consumer) Listen() {
	fmt.Println("Consumer is running!")
	pubsub := c.client.Subscribe(c.channel)

	_, err := pubsub.Receive()
	if err != nil {
		fmt.Println(err)
		return
	}

	ch := pubsub.Channel()
	for msg := range ch {
		fmt.Println(msg.Channel, msg.Payload)
		message := Message{}
		bytesMessage := []byte(msg.Payload)
		err = json.Unmarshal(bytesMessage, &message)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%+v\n", message)
		connections := c.connectionMap.GetConnectionsByThreadId(message.ThreadId)
		for _, c := range connections {
			err := wsutil.WriteServerMessage(c.conn, ws.OpText, bytesMessage)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	fmt.Println("Consumer done")
}

func NewConsumer(client *redis.Client, connectionMap *ConnectionMap, channel string) Consumer {
	pubsub := client.Subscribe(channel)
	return Consumer{
		client:        client,
		pubsub:        pubsub,
		connectionMap: connectionMap,
		channel:       channel,
	}
}
