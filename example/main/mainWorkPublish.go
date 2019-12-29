package main

import (
	"RabbitMQ"
	"fmt"
	"strconv"
)

// Work Mode is for load balance
//Every message will be consumed only once
func main() {
	rabbitmq := RabbitMQ.NewRabbitMQSimple("queueWork")
	for i := 1; i <= 100; i++ {
		rabbitmq.PublishSimple("Message Num" + strconv.Itoa(i))
		fmt.Printf("send message %v success\n", i)
	}
}
