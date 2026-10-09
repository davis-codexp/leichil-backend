package message_broker

import (
	"context"
	"encoding/json"
	"log"
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

func InitPublisher() {
	var conn *amqp.Connection = GetBrokerConnection()
	if conn == nil {
		return
	}
	var channel *amqp.Channel = GetChannel(conn, queue_name)
	if channel == nil {
		return
	}
	publish_channel = channel
}

func SendSms(message_id, mobile, text string) {
	if publish_channel == nil {
		log.Print("Publish Channel is not initialized.")
		return
	}
	var body = SmsBody{
		Mobile:    mobile,
		MessageID: message_id,
		Text:      text,
	}
	PushToQueue(publish_channel, queue_name, body)
}

func PushToQueue(channel *amqp.Channel, queue_name string, body SmsBody) bool {
	data, err := json.Marshal(body)
	if err != nil {
		log.Print("Error in JSON Marshal of sms body: ", err.Error())
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
