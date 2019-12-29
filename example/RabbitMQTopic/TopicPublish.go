package main

import (
	"RabbitMQ"
	"fmt"
	"strconv"
	"time"
)

func main() {
	rabbitmq1 := RabbitMQ.NewRabbitMQTopic("topicExchange", "TopicScience.TopicPhysics")
	rabbitmq2 := RabbitMQ.NewRabbitMQTopic("topicExchange", "TopicScience.TopicPhysics.TopicQuantum")

	for i := 0; i <= 100; i++ {
		rabbitmq1.PublishTopic("Publish message Routing to TopicPhysics, Num: " + strconv.Itoa(i))
		fmt.Printf("Message  %v send to Topic Physics Success\n", i)
		rabbitmq2.PublishTopic("Publish message Routing to TopicQuantum Num: " + strconv.Itoa(i))
		fmt.Printf("Message  %v send to Topic Quantum Success\n", i)
		time.Sleep(time.Second)
	}
}
