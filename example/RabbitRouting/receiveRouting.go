package main

import "RabbitMQ"

func main(){
	rabbitmq := RabbitMQ.NewRabbitMQRouteing("routingExchange","routeKeyOne")
	rabbitmq.ReceiveRouting()
}
