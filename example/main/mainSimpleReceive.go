package main

import(
	"RabbitMQ"
	"fmt"
)
func main(){
	rabbitmq := RabbitMQ.NewRabbitMQSimple("queueSimple")
	rabbitmq.ConsumeSimple()
	fmt.Println("message receive success")
}
