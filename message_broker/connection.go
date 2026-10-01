package message_broker

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func GetBrokerConnection() *amqp.Connection {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	username := os.Getenv("BROKER_USERNAME")
	password := os.Getenv("BROKER_PASSWORD")
	address := os.Getenv("BROKER_ADDRESS")
	port := os.Getenv("BROKER_PORT")
	broker := fmt.Sprintf("amqp://%s:%s@%s:%s", username, password, address, port)
	conn, err := amqp.Dial(broker)
	if err != nil {
		panic(err)
	}
	return conn
}

func GetChannel(conn *amqp.Connection, queue_name string) *amqp.Channel {
	channel, err := conn.Channel()

	if err != nil {
		panic(err)
	}

	channel.QueueDeclare(queue_name, false, false, false, false, nil)
	return channel
}
