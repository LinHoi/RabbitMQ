package main

import (
	"RabbitMQ"
	"fmt"
	"strconv"
	"time"
)

func main(){
	rabbitmq := RabbitMQ.NewRabbitMQPubSub("TestExchange")
	for i := 0;i <= 100; i++ {
		rabbitmq.PublishPub("Publsih Message Num:"+strconv.Itoa(i))
		fmt.Printf("Publish message Num: %v send \n",i)
		time.Sleep(1*time.Second)
	}
}
