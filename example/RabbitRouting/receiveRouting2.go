package main

import "RabbitMQ"

func main(){
	rabbitmq := RabbitMQ.NewRabbitMQRouteing("routingExchange","routeKeyTwo")
	rabbitmq.ReceiveRouting()
}

