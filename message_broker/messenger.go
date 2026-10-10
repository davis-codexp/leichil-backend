package message_broker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SmsBody struct {
	Mobile    string `json:"mobile"`
	MessageID string `json:"message_id"`
	Text      string `json:"text"`
}

var publish_channel *amqp.Channel

const queue_name = "sms_queue"

var (
	instance_conn *amqp.Connection
	conn_mutex    sync.RWMutex
)

func InitPublisher() {
	var conn *amqp.Connection = get_broker_connection()
	if conn == nil {
		return
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Println("Error in creating channel: ", err.Error())
		return
	}
	defer ch.Close()
	_, err = ch.QueueDeclare(queue_name, false, false, false, false, nil)
	if err != nil {
		log.Println("Error declaring queue: ", err.Error())
	}
}

func get_broker_connection() *amqp.Connection {
	conn_mutex.RLock()
	if instance_conn != nil && !instance_conn.IsClosed() {
		defer conn_mutex.RUnlock()
		return instance_conn
	}
	conn_mutex.RUnlock()
	conn_mutex.Lock()
	defer conn_mutex.Unlock()
	if instance_conn != nil && !instance_conn.IsClosed() {
		return instance_conn
	}
	var (
		username = os.Getenv("BROKER_USERNAME")
		password = os.Getenv("BROKER_PASSWORD")
		address  = os.Getenv("BROKER_ADDRESS")
		port     = os.Getenv("BROKER_PORT")
		broker   = fmt.Sprintf("amqp://%s:%s@%s:%s", username, password, address, port)
	)
	conn, err := amqp.Dial(broker)
	if err != nil {
		fmt.Println("Error in connecting to Rabbit MQ: ", err.Error())
		return nil
	}
	instance_conn = conn

	//Background Monitor to reconnect to Message Broker in case of a connection loss
	go func() {
		var close_chan chan *amqp.Error = make(chan *amqp.Error, 1)
		conn.NotifyClose(close_chan)
		err_reason := <-close_chan
		log.Println("Connection to Message Broker Lost. Reconnecting... ", err_reason)

		for {
			time.Sleep(5 * time.Second)
			if next_conn := get_broker_connection(); next_conn != nil {
				log.Println("Reconnected to Message Broker")
				break
			}
		}
	}()

	return instance_conn
}

func SendSms(message_id, mobile, text string) {
	var conn *amqp.Connection = get_broker_connection()
	if conn == nil {
		log.Println("Unable to connecto to Message Broker")
		return
	}
	channel, err := conn.Channel()
	if err != nil {
		log.Println("Unable to get channel: ", err.Error())
		return
	}
	defer channel.Close()

	var body = SmsBody{
		Mobile:    mobile,
		MessageID: message_id,
		Text:      text,
	}
	PushToQueue(channel, queue_name, body)
}

func PushToQueue(channel *amqp.Channel, queue_name string, body SmsBody) bool {
	data, err := json.Marshal(body)
	if err != nil {
		log.Println("Error in JSON Marshal of sms body: ", err.Error())
		return false
	}
	var message = amqp.Publishing{
		ContentType: "application/json",
		Body:        data,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = channel.PublishWithContext(ctx, "", queue_name, false, false, message)
	if err != nil {
		log.Print("Unable to push SMS payload to Message Queue: ", err.Error())
		return false
	}
	return true
}
