package main

import (
	"RabbitMQ"
	"fmt"
)
func main(){
	rabbitmq := RabbitMQ.NewRabbitMQSimple("queueSimple")
	rabbitmq.PublishSimple("Hello RabbitMQ")
	fmt.Println("send message success")
}

