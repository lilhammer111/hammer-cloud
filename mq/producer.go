package mq

import (
	"github.com/lilhammer111/hammer-cloud/config"
	"github.com/streadway/amqp"
	"log"
)

var (
	conn    *amqp.Connection
	channel *amqp.Channel
)

func init() {
	if !initChannel() {
		log.Fatal("Failed to initialize the channel")
	}
}

func initChannel() bool {

	// 1. judge if the channel had been created
	if channel != nil {
		return true
	}
	// 2. get a rabbitmq connection
	var err error
	conn, err = amqp.Dial(config.RabbitURL)
	if err != nil {
		log.Println("dial err")
		return false
	}
	// 3. open a channel for consume or publish message
	channel, err = conn.Channel()
	if err != nil {
		log.Println("create channel err")
		closeConnection()
		return false
	}
	return true
}

func Publish(exchange, routingKey string, msg []byte) bool {
	// 1. judge if channel is working
	if !initChannel() {
		return false
	}
	// 2. publish msg
	err := channel.Publish(exchange, routingKey, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        msg,
	})
	if err != nil {
		closeConnection()
		return false
	}
	return true
}

func closeConnection() {
	if channel != nil {
		channel.Close()
		channel = nil
	}
	if conn != nil {
		conn.Close()
		conn = nil
	}
}
