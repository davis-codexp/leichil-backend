package message_broker

import (
	"fmt"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func GetBrokerConnection() *amqp.Connection {
	username := os.Getenv("BROKER_USERNAME")
	password := os.Getenv("BROKER_PASSWORD")
	address := os.Getenv("BROKER_ADDRESS")
	port := os.Getenv("BROKER_PORT")
	broker := fmt.Sprintf("amqp://%s:%s@%s:%s", username, password, address, port)
	conn, err := amqp.Dial(broker)
	if err != nil {
		fmt.Println("Error in connecting to Rabbit MQ: ", err.Error())
		return nil
	}
	return conn
}

func GetChannel(conn *amqp.Connection, queue_name string) *amqp.Channel {
	channel, err := conn.Channel()
	if err != nil {
		fmt.Println("Error in initializing Rabbit MQ channel: ", err.Error())
		return nil
	}
	_, err = channel.QueueDeclare(queue_name, false, false, false, false, nil)
	if err != nil {
		fmt.Println("Error in declaring Message Queue: ", err.Error())
		channel.Close()
		return nil
	}
	return channel
}
