package main

import (
	"RabbitMQ"
	"fmt"
	"strconv"
	"time"
)

func main(){
	rabbitmq1 := RabbitMQ.NewRabbitMQRouteing("routingExchange","routeKeyOne")
	rabbitmq2 := RabbitMQ.NewRabbitMQRouteing("routingExchange","routeKeyTwo")

	for i := 0; i <= 100; i++ {
		rabbitmq1.PublishRouting("Publish message Routing to RouteOne, Num: "+ strconv.Itoa(i))
		fmt.Printf("Message  %v send to Route One Success\n", i)
		rabbitmq2.PublishRouting("Publish message Routing to RouteTwo, Num: "+ strconv.Itoa(i))
		fmt.Printf("Message  %v send to Route Two Success\n", i)
		time.Sleep(time.Second)
	}
}
