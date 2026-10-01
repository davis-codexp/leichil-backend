package message_broker

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SmsBody struct {
	Mobile    string `json:"mobile"`
	MessageID string `json:"message_id"`
	Text      string `json:"text"`
}

func SendSms(message_id, mobile, text string) {
	const queue_name = "sms_queue"
	conn := GetBrokerConnection()
	if conn == nil {
		return
	}
	defer conn.Close()
	channel := GetChannel(conn, queue_name)
	if channel == nil {
		return
	}
	body := SmsBody{
		Mobile:    mobile,
		MessageID: message_id,
		Text:      text,
	}
	PushToQueue(channel, queue_name, body)
}

func PushToQueue(channel *amqp.Channel, queue_name string, body SmsBody) bool {
	data, err := json.Marshal(body)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	message := amqp.Publishing{
		ContentType: "text/plain",
		Body:        data,
	}
	err = channel.PublishWithContext(context.TODO(), "", queue_name, false, false, message)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}
