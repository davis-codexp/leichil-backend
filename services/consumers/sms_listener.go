package consumers

import (
	"bytes"
	"encoding/json"
	"leichil-backend/message_broker"
	"leichil-backend/services/models"
	"net/http"
	"net/url"
	"os"

	"time"

	"github.com/gofiber/fiber/v3/log"
)

var http_client = &http.Client{
	Timeout: 15 * time.Second,
}

type SmsConfig struct {
	URL      string
	AuthKey  string
	Route    string
	SenderID string
}

func SmsListener() {
	const queue_name = "sms_queue"
	var config = SmsConfig{
		URL:      os.Getenv("SMS_URL"),
		AuthKey:  os.Getenv("SMS_AUTH"),
		Route:    os.Getenv("SMS_ROUTE"),
		SenderID: os.Getenv("SENDER_ID"),
	}
	for {
		var conn = message_broker.GetBrokerConnection()
		if conn == nil {
			log.Error("Error in connecting to Rabbit MQ")
			time.Sleep(5 * time.Second)
			continue
		}
		var ch = message_broker.GetChannel(conn, queue_name)
		if ch == nil {
			log.Error("Error in fetching channel: ")
			conn.Close()
			time.Sleep(5 * time.Second)
			continue
		}

		messages, err := ch.Consume(queue_name, "", true, false, false, false, nil)
		if err != nil {
			log.Error("Failed to register consumer: ", err.Error())
			ch.Close()
			conn.Close()
			time.Sleep(5 * time.Second)
			continue
		}
		log.Info("SMS Listener Started")
		for data := range messages {
			var body = map[string]string{}
			if err := json.Unmarshal(data.Body, &body); err != nil {
				log.Error("Malformed message queue delivery package: ", err.Error())
				continue
			}
			var message_id = url.QueryEscape(body["message_id"])
			var mobile = body["mobile"]
			var text = body["text"]
			send_sms(mobile, message_id, text, config)
		}
		log.Warn("RabbitMQ channel connection broker. Reconnecting...")
		ch.Close()
		conn.Close()
		time.Sleep(2 * time.Second)
	}
}

func send_sms(mobile, message_id, text string, config SmsConfig) {
	var body = models.SmsBody{
		Route:          config.Route,
		SenderID:       config.SenderID,
		Message:        message_id,
		VariableValues: text,
		Numbers:        mobile,
	}
	json_data, err := json.Marshal(body)
	if err != nil {
		log.Error("Error in JSON Marshal of SMS Body: ", err.Error())
		return
	}
	req, err := http.NewRequest("POST", config.URL, bytes.NewBuffer(json_data))
	if err != nil {
		log.Error("Error in creating SMS body: ", err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", config.AuthKey)
	response, err := http_client.Do(req)
	if err != nil {
		log.Error("Error in seding SMS: ", err.Error())
		return
	}
	defer response.Body.Close()
	log.Info("SMS API Response Status: ", response.Status)
}
