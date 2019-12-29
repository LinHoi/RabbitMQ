package main

import (
	"RabbitMQ"
	"fmt"
)
func main(){
	rabbitmq := RabbitMQ.NewRabbitMQSimple("queueWork")
	rabbitmq.ConsumeSimple()
	fmt.Println("receive message success")
}
