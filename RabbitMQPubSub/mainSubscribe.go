package main

import "RabbitMQ"

func main(){
	rabbitmq := RabbitMQ.NewRabbitMQPubSub("TestExchange")
	rabbitmq.ReceiveSub()

}
