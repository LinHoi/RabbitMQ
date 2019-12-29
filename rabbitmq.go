package RabbitMQ

import (
	"fmt"
	"github.com/go-acme/lego/log"
	"github.com/streadway/amqp"
)

//url :   amqp://username:password@rabbitserver_address:server_port/virtual_host
const MQURL = "amqp://hikari:123456@127.0.0.1:5672/host"

type RabbitMQ struct {
	conn *amqp.Connection
	channel *amqp.Channel

	QueueName 	string

	Exchange 	string

	Key			string

	Mqurl		string
}

func NewRabbitMQ(queueName, exchange, key string) *RabbitMQ {
	rabbitmq := &RabbitMQ{
		QueueName: queueName,
		Exchange:  exchange,
		Key:       key,
		Mqurl:     MQURL,
	}

	var err error
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "connect to RabbitMQ fail")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "getting channel fail")
	return  rabbitmq
}


func (r *RabbitMQ) Destory () {
	r.channel.Close()
	r.conn.Close()
}

func (r *RabbitMQ) failOnErr(err error, message  string){
	if err != nil {
		log.Fatal("%s:%s",message,err)
		panic(fmt.Sprintf("%s:%s",message,err))
	}
}

//========================================================
//Simple Mode
//========================================================
//Simple Mode Step 1 : generate a simple instance of RabbitMQ
func NewRabbitMQSimple(queueName string) *RabbitMQ{
	rabbitmq := NewRabbitMQ(queueName,"","")
	var err error
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "connect to RabbitMQ fail")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "getting channel fail")
	return  rabbitmq
}

//Simple Mode Step 2 : Product
func (r *RabbitMQ) PublishSimple(message string){
	//1. apply for a queue
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		false,
		false,
		false,
		false,
		nil,
		)
	if err != nil {
		fmt.Println(err)
	}
	//2. send message to queue
	err = r.channel.Publish(
		r.Exchange,
		r.QueueName,
		//if mandatory == true, it will send back the message to sender when the queue can't be found by exchange and key
		false,
		//if immediate == true,it will send back the message when the consumer is not found
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})

	if err != nil {
		panic(err)
	}
}

//Simple Mode Step 3: Consumer
func (r *RabbitMQ) ConsumeSimple(){
	//1. apply for a queue
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}

	//2. receive message
	msgs, err := r.channel.Consume(
		r.QueueName,
		//saperate different consumer
		"",
		//ack to RabbitMQ when queue is consumed
		true,
		false,
		//if noLocal ==true, message can't be sent to this connection itself
		false,
		false,
		nil,
		)
	if err != nil {
		fmt.Println(err)
	}

	//3. consume queue
	forever := make (chan bool)
	go func(){
		for d := range msgs{
			log.Printf("Receive a message: %s",d.Body)
		}
	}()
	log.Printf("Waiting for messages, press CTRL + S to exit")
	<- forever

}

//=======================================================
// Publish and Subscribe Mode
//======================================================
//Publish/Subscribe Mode
func NewRabbitMQPubSub(exchange string) *RabbitMQ{
	rabbitmq := NewRabbitMQ("",exchange,"")
	var err error

	rabbitmq.conn, err = amqp.Dial(MQURL)
	rabbitmq.failOnErr(err,"fail to connect to rabbitmq server")

	rabbitmq.channel,err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "fail to open a channel")
	return rabbitmq
}

//Pub/Sub Mode : Publish
func (r *RabbitMQ) PublishPub(message string){
	//1 . generate a exchange
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
		)

	r.failOnErr(err, "Fail to declare a exchange")

	//2. publish/send message
	err = r.channel.Publish(
		r.Exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:		[]byte(message),
		})
}

//Pub/Sub Mode : receive
func (r *RabbitMQ) ReceiveSub(){
	//1 . generate a exchange
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		//the kind of exchange
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Fail to declare a exchange")

	//2. generate queue with no queue name
	q, err := r.channel.QueueDeclare(
		"",//queue name is random
		false,
		false,
		true,
		false,
		nil,
		)
	r.failOnErr(err, "Fail to declare a queue")

	//3. bind queue with exchange
	err = r.channel.QueueBind(
		q.Name,
		"",
		r.Exchange,
		false,
		nil,
		)
	r.failOnErr(err," Fail to bind queue with exchange")

	//4. consume
	message, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
		)

	forever := make(chan bool)
	go func(){
		for d := range message {
			log.Printf("Received a message: %s", d.Body)
		}
	}()
	fmt.Println("Press CTRL + C to quit")
	<- forever
}



//================================================
// Route Mode
//================================================
// Route Mode is an upgrade Version of Pub/Sub Mode
// Route can assign the queue name that can receive the message
// It use routingKey and direct kind of exchange to achieve this function
// The Code of Route Mode and Pub/Sub are similar
func NewRabbitMQRouteing(exchange, routingKey string) *RabbitMQ{
	rabbitmq := NewRabbitMQ("",exchange,routingKey)
	var err error

	rabbitmq.conn, err = amqp.Dial(MQURL)
	rabbitmq.failOnErr(err,"fail to connect to rabbitmq server")

	rabbitmq.channel,err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "fail to open a channel")
	return rabbitmq
}

//Route Mode : Publish
func (r *RabbitMQ) PublishRouting(message string){
	//1 . generate a exchange
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)

	r.failOnErr(err, "Fail to declare a exchange")

	//2. send message
	err = r.channel.Publish(
		r.Exchange,
		r.Key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:		[]byte(message),
		})
}

//Route Mode : receive
func (r *RabbitMQ) ReceiveRouting(){
	//1 . generate a exchange
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		//the kind of exchange
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Fail to declare a exchange")

	//2. generate queue with no queue name
	q, err := r.channel.QueueDeclare(
		"",//queue name is random
		false,
		false,
		true,
		false,
		nil,
	)
	r.failOnErr(err, "Fail to declare a queue")

	//3. bind queue with exchange
	err = r.channel.QueueBind(
		q.Name,
		r.Key,
		r.Exchange,
		false,
		nil,
	)
	r.failOnErr(err," Fail to bind queue with exchange")

	//4. consume
	message, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	forever := make(chan bool)
	go func(){
		for d := range message {
			log.Printf("Received a message: %s", d.Body)
		}
	}()
	fmt.Println("Press CTRL + C to quit")
	<- forever
}

//================================================
// Topic Mode
//================================================
// Route Mode is an upgrade Version of Route Mode
// Topic Mode can assign the queue_name regular expression that can receive the message
// It use topic kind of exchange to achieve this function
// The Code of Topic Mode and Route Mode are all most the same
// In routingKey, "*" matches a word, "#" matches 0-some words, every word should be separate by "."
func NewRabbitMQTopic(exchange, routingKey string) *RabbitMQ{
	rabbitmq := NewRabbitMQ("",exchange,routingKey)
	var err error

	rabbitmq.conn, err = amqp.Dial(MQURL)
	rabbitmq.failOnErr(err,"fail to connect to rabbitmq server")

	rabbitmq.channel,err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "fail to open a channel")
	return rabbitmq
}

//Route Mode : Publish
func (r *RabbitMQ) PublishTopic(message string){
	//1 . generate a exchange
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)

	r.failOnErr(err, "Fail to declare a exchange")

	//2. send message
	err = r.channel.Publish(
		r.Exchange,
		r.Key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:		[]byte(message),
		})
}

//Route Mode : receive
func (r *RabbitMQ) ReceiveTopic(){
	//1 . generate a exchange
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		//the kind of exchange
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Fail to declare a exchange")

	//2. generate queue with no queue name
	q, err := r.channel.QueueDeclare(
		"",//queue name is random
		false,
		false,
		true,
		false,
		nil,
	)
	r.failOnErr(err, "Fail to declare a queue")

	//3. bind queue with exchange
	err = r.channel.QueueBind(
		q.Name,
		r.Key,
		r.Exchange,
		false,
		nil,
	)
	r.failOnErr(err," Fail to bind queue with exchange")

	//4. consume
	message, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	forever := make(chan bool)
	go func(){
		for d := range message {
			log.Printf("Received a message: %s", d.Body)
		}
	}()
	fmt.Println("Press CTRL + C to quit")
	<- forever
}