package main

import "RabbitMQ"

func main(){
	rabbitmq := RabbitMQ.NewRabbitMQTopic("topicExchange","#")
	rabbitmq.ReceiveTopic()

}
