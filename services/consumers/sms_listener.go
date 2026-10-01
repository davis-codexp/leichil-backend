package consumers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"leichil-backend/message_broker"
	"leichil-backend/services/models"
	"net/http"
	"net/url"
	"os"

	"github.com/gofiber/fiber/v3/log"
	"github.com/joho/godotenv"
)

func SmsListener() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	const queue_name = "sms_queue"
	conn := message_broker.GetBrokerConnection()
	defer conn.Close()
	ch := message_broker.GetChannel(conn, queue_name)
	defer ch.Close()
	messages, err := ch.Consume(queue_name, "", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}
	forever := make(chan bool)
	var body = map[string]string{}
	go func() {
		for data := range messages {
			json.Unmarshal(data.Body, &body)
			var message_id string = url.QueryEscape(body["message_id"])
			var mobile string = body["mobile"]
			var text string = body["text"]
			send_sms(mobile, message_id, text)
		}
	}()
	fmt.Println("SMS Listener Started")
	<-forever
}

func send_sms(mobile, message_id, text string) {
	var url = os.Getenv("SMS_URL")
	var auth_key = os.Getenv("SMS_AUTH")
	var body = models.SmsBody{
		Route:          os.Getenv("SMS_ROUTE"),
		SenderID:       os.Getenv("SENDER_ID"),
		Message:        message_id,
		VariableValues: text,
		Numbers:        mobile,
	}
	json_data, err := json.Marshal(body)
	if err != nil {
		log.Error("Error in JSON Marshal of SMS Body: ", err)
		return
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json_data))
	if err != nil {
		log.Error("Error in SMS Request Creation: ", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", auth_key)
	var client *http.Client = &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Error("Error in SMS API Call: ", err)
		return
	}
	defer response.Body.Close()
	log.Info("SMS API Response Status: ", response.Status)
}
