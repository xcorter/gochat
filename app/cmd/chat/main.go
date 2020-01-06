package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/xcorter/gochat/internal/chat"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	fmt.Println("Start")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
	})

	connMap := chat.NewConnectionMap()
	redisClient := NewRedis()
	channel := "chat"
	publisher := chat.NewPublisher(redisClient, channel)
	consumer := chat.NewConsumer(redisClient, &connMap, channel)
	go consumer.Listen()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		threadIdArr, ok := r.URL.Query()["threadId"]
		if !ok {
			fmt.Println(w, "bad request")
			return
		}

		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			fmt.Println(err)
			return
		}

		threadIdStr := threadIdArr[0]
		i, err := strconv.Atoi(threadIdStr)
		if err != nil {
			fmt.Println(err)
			return
		}
		connMap.NewConnection(conn, i)

		go func() {
			log.Printf("new request")
			for {
				msg, _, err := wsutil.ReadClientData(conn)
				if err != nil {
					fmt.Println(err)
					return
				}

				message := chat.Message{}
				err = json.Unmarshal(msg, &message)
				if err != nil {
					fmt.Println(err)
					return
				}
				publisher.Publish(message)
			}
		}()
	})

	port := ":" + os.Getenv("API_PORT")
	http.ListenAndServe(port, nil)
}

func NewRedis() *redis.Client {
	redisConnection := os.Getenv("REDIS_CONNECTION")
	client := redis.NewClient(&redis.Options{
		Addr:     redisConnection,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return client
}
